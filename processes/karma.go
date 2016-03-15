package processes

import (
	"encoding/json"
	"errors"

	"menteslibres.net/gosexy/to"

	"github.com/robxu9/kahinah/models"
)

const (
	// Old behaviours:

	// KarmaLowHigh specifies that if the total amount of karma is less than
	// or equal to the lowest possible negative value, submission can be
	// rejected whilst the total amount of karma greater than or equal to the
	// highest possible positive value enables submission. Values in between
	// promote no action.
	KarmaLowHigh = "LowHigh"
	// KarmaLowAny specifies that if the total amount of karma is less than
	// or equal to the lowest possible negative value, submission can be
	// rejected whilst any other value enables submission.
	KarmaLowAny = "LowAny"
	// KarmaAnyHigh specifies that the total amount of karma is greater than
	// or equal to the highest possible positive value enables submission
	// whilst any other value allows the submission to be rejected.
	KarmaAnyHigh = "AnyHigh"

	// New behaviours (Karma):

	// KarmaMaxWithBlock specifies that:
	// The lowest possible negative value, if present, blocks a submit, while
	// the highest possible positive value is required to enable submit. There
	// must be at least one positive value, or else submit will never be
	// enabled. To permit blocking submits, ensure a negative value is defined.
	KarmaMaxWithBlock = "MaxWithBlock"
	// KarmaAnyWithBlock specifies that:
	// The lowest possible negative value, if present, blocks a submit, Any
	// other value enables a submit. To permit blocking submits, ensure that a
	// negative value is defined.
	KarmaAnyWithBlock = "AnyWithBlock"
	// KarmaMaxNoBlock specifies that:
	// The highest possible positive value is required to enable submit, but the
	// lowest possible negative value will not block the change.
	KarmaMaxNoBlock = "MaxNoBlock"
	// KarmaNoOp specifies that:
	// The label is purely informational and values are not considered when
	// determining whether a change is submittable.
	KarmaNoOp = "NoOp"
	// KarmaLockList specifies that:
	// This label does not block nor submit. Instead, it prevents the karma
	// process from ever succeeding or failing, and therefore the stage from
	// ever progressing forward. Valid values are 0 (allowed to move forward)
	// or 1 (prohibited from moving forward).
	KarmaLockList = "LockList"
)

const (
	KarmaAccept = "accept"
	KarmaReject = "reject"

	// KarmaStateWaiting - waiting for input
	KarmaStateWaiting int64 = iota
	KarmaStateAccepted
	KarmaStateRejected

	karmaCanVote = 1 << iota
	karmaCanAccept
	karmaCanReject
)

var (
	ErrKarmaLabelNotFound    = errors.New("process/karma: label not found")
	ErrKarmaLabelIsClosed    = errors.New("process/karma: can't take action on that label")
	ErrKarmaLabelCannotClose = errors.New("process/karma: can't close the label")
	ErrKarmaValueOutOfRange  = errors.New("process/karma: value out of range (check metadata?)")

	karmaGenericOk = map[string]interface{}{
		"success": true,
	}
)

func newKarmaProcess(m *models.ListStageProcess, cfg map[string]interface{}) Process {
	k := &Karma{
		model:  m,
		Labels: []*KarmaLabel{},
	}

	// unmarshal data
	if m.Data != nil {
		json.Unmarshal(m.Data, k)
	} else {
		k.State = map[string]int64{}
		k.Votes = map[string]map[string]int64{}
	}

	// and set configuration, which in yaml looks like this:
	// karma:
	//   enable: true
	//   labels:
	//      Code-Review:
	//         function: LowHigh
	//         low: -3
	//         high: 3
	labels := cfg["labels"].(map[string]interface{})
	for label, opts := range labels {
		optsCast := opts.(map[string]interface{})
		k.Labels = append(k.Labels, &KarmaLabel{
			Label:     label,
			Function:  optsCast["function"].(string),
			LowValue:  optsCast["low"].(int64),
			HighValue: optsCast["high"].(int64),
		})
	}

	return k
}

func init() {
	Mapping["karma"] = newKarmaProcess
}

// Karma represents the data for the current karma checking process. We contain
// both the old model, where a consensus was considered acceptable, or the
// the new model, which is based on Gerrit:
// https://gerrit-review.googlesource.com/Documentation/config-labels.html#label_function
type Karma struct {
	model  *models.ListStageProcess
	Labels []*KarmaLabel `json:"-"`

	Votes map[string]map[string]int64 // label ~> username ~> value
	State map[string]int64            // label ~> state
}

type KarmaLabel struct {
	Label     string
	Function  string
	LowValue  int64
	HighValue int64
}

func (k *Karma) Start() error {
	// we don't need to start anything
	return nil
}

func (k *Karma) Stop() error {
	// we don't exactly need to stop anything
	return nil
}

func (k *Karma) Reset() error {
	// remove all previous karma.
	k.Votes = map[string]map[string]int64{}
	k.State = map[string]int64{}
	k.save()
	return nil
}

func (k *Karma) save() {
	bte, err := json.Marshal(k)
	if err != nil {
		panic(err)
	}

	k.model.Data = bte
}

func (k *Karma) Finalise() error {
	// we don't need to finalise anything
	return nil
}

func (k *Karma) Status() ProcessStatus {
	// return the state
	for _, v := range k.Labels {
		// special cases:
		if v.Label == KarmaNoOp {
			continue // doesn't count
		}

		s, ok := k.State[v.Label]
		if !ok {
			return ProcessRunning
		}
		if s == KarmaStateWaiting {
			return ProcessRunning
		} else if s == KarmaStateRejected {
			return ProcessFail
		}
	}
	return ProcessOK
}

// old behaviour: calculate the total amount of karma on a label
func (k *Karma) calculateTotal(kl *KarmaLabel) int64 {
	votes := k.Votes[kl.Label]
	if votes == nil {
		return 0 // welp, guess we got nothing
	}

	var total int64
	for _, v := range votes {
		total += v
	}

	return total
}

// new behaviour: get lowest value and highest value
// returns (lowest, highest)
func (k *Karma) getLowestHighest(kl *KarmaLabel) (int64, int64) {
	votes := k.Votes[kl.Label]
	if votes == nil {
		return 0, 0 // welp..
	}

	var highest int64
	var lowest int64
	for _, v := range votes {
		if v > highest {
			highest = v
		}
		if v < lowest {
			lowest = v
		}
	}

	return lowest, highest
}

// checkLabel looks for a label, and then indicates whether we can vote on
// it, as well as accept/reject (or neither)
func (k *Karma) checkLabel(kl *KarmaLabel) int {
	// check if it's actually already been accepted/rejected
	if k.State[kl.Label] != KarmaStateWaiting {
		return 0 // we're not waiting, so we can't vote/accept/reject
	}

	result := karmaCanVote // not accepted/rejected? we can at least vote then

	// you must be able to at least vote in order to do anything
	switch kl.Function {
	case KarmaLowHigh:
		total := k.calculateTotal(kl)
		if total <= kl.LowValue {
			result = result | karmaCanReject
		}
		if total >= kl.HighValue {
			result = result | karmaCanAccept
		}

	case KarmaLowAny:
		total := k.calculateTotal(kl)
		if total <= kl.LowValue {
			result = result | karmaCanReject
		} else {
			result = result | karmaCanAccept
		}

	case KarmaAnyHigh:
		total := k.calculateTotal(kl)
		if total >= kl.HighValue {
			result = result | karmaCanAccept
		} else {
			result = result | karmaCanReject
		}

	case KarmaMaxWithBlock:
		low, high := k.getLowestHighest(kl)
		if high > 0 && high >= kl.HighValue {
			result = result | karmaCanAccept
		}

		if low < 0 && low <= kl.LowValue {
			result = result | karmaCanReject
		}

	case KarmaAnyWithBlock:
		low, _ := k.getLowestHighest(kl)
		if low < 0 && low <= kl.LowValue {
			result = result | karmaCanReject
		} else {
			result = result | karmaCanAccept
		}
	case KarmaMaxNoBlock:
		_, high := k.getLowestHighest(kl)
		result = result | karmaCanReject
		if high >= kl.HighValue {
			result = result | karmaCanAccept
		}
	case KarmaNoOp:
		result = 0
	case KarmaLockList:
		result = result | karmaCanAccept
	}

	return result
}

func (k *Karma) APIRequest(user *models.User, action string, value string) (interface{}, error) {
	// action == label
	// value == -1/0/+1/etc

	// check if we have the action available
	for _, v := range k.Labels {
		if v.Label != action {
			continue
		}

		// check if we can take action on this label
		votable := k.checkLabel(v)
		if votable&karmaCanVote == 0 {
			return nil, ErrKarmaLabelIsClosed
		}

		if value == KarmaAccept {
			if votable&karmaCanAccept == 0 {
				return nil, ErrKarmaLabelCannotClose
			}
			k.State[v.Label] = KarmaStateAccepted

			return karmaGenericOk, nil
		} else if value == KarmaReject {
			if votable&karmaCanReject == 0 {
				return nil, ErrKarmaLabelCannotClose
			}
			k.State[v.Label] = KarmaStateRejected

			return karmaGenericOk, nil
		}

		vote := to.Int64(value)
		if vote < v.LowValue || vote > v.HighValue {
			return nil, ErrKarmaValueOutOfRange
		}

		// alright, add it in
		labelVotes, ok := k.Votes[v.Label]
		if !ok {
			labelVotes = map[string]int64{}
			k.Votes[v.Label] = labelVotes
		}

		labelVotes[user.Username] = vote

		k.save()

		return karmaGenericOk, nil
	}

	return nil, ErrKarmaLabelNotFound
}

func (k *Karma) APIMetadata(user *models.User) interface{} {
	// check every single label
	votable := map[string]int{}
	for _, v := range k.Labels {
		votable[v.Label] = k.checkLabel(v)
	}

	// return list of karma + adding karma
	return map[string]interface{}{
		"votes":   k.Votes,
		"state":   k.State,
		"labels":  k.Labels,
		"votable": votable,
	}
}
