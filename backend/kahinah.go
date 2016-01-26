package kahinah

import rt "gopkg.in/dancannon/gorethink.v1"

// K represents the main entry point for Kahinah operations, including direct
// access to the database. All operations are thread-safe.
type K struct {
	Session *rt.Session
	Opts    *KOpts
	Rules   []*Rule
}

// KOpts represents options to opening the main K instance.
type KOpts struct {
	DBConnectOpts rt.ConnectOpts
	DBRunOpts     rt.RunOpts
	DBName        string
	DBInit        bool
}

// WRError represents an error emitted from RethinkDB's response.
type WRError struct {
	FirstError string
}

// Error is the error emitted.
func (w *WRError) Error() string {
	return w.FirstError
}

// Open opens a connection to the K instance. If there is a database error, it
// will return an error of type DBError.
func Open(opts *KOpts) (*K, error) {
	session, err := rt.Connect(opts.DBConnectOpts)
	if err != nil {
		return nil, &DBError{err}
	}

	// creates the db. If the db already exists, we don't do anything
	if _, err = rt.DBCreate(opts.DBName).RunWrite(session, rt.RunOpts{}); err == nil {
		// db didn't exist before, create tables & indexes
		if _, err = rt.DB(opts.DBName).TableCreate(TableAdvisory, rt.TableCreateOpts{}).RunWrite(session, opts.DBRunOpts); err != nil {
			return nil, &DBError{err}
		}
		if _, err = rt.DB(opts.DBName).TableCreate(TableCounter, rt.TableCreateOpts{
			PrimaryKey: "group",
		}).RunWrite(session, opts.DBRunOpts); err != nil {
			return nil, &DBError{err}
		}
		if _, err = rt.DB(opts.DBName).TableCreate(TableUpdate, rt.TableCreateOpts{}).RunWrite(session, opts.DBRunOpts); err != nil {
			return nil, &DBError{err}
		}

		// see advisory.go for the Advisory struct and the parts that make up
		// this compound index.
		if _, err = rt.DB(opts.DBName).Table(TableAdvisory).IndexCreateFunc("advisory_id", func(row rt.Term) interface{} {
			return []interface{}{row.Field("group"), row.Field("year"), row.Field("advisory_num")}
		}).RunWrite(session, opts.DBRunOpts); err != nil {
			return nil, &DBError{err}
		}

		// see update.go for the Update struct and the parts that make up these
		// single indexes
		indexNames := []string{"platform", "name", "evr", "submitter", "created_at", "type", "connector"}
		for _, v := range indexNames {
			if _, err = rt.DB(opts.DBName).Table(TableUpdate).IndexCreate(v).RunWrite(session, opts.DBRunOpts); err != nil {
				return nil, &DBError{err}
			}
		}
	}

	return &K{
		Session: session,
		Opts:    opts,
	}, nil
}

// Close closes the Kahinah database.
func (k *K) Close() error {
	return k.Session.Close()
}

// DB starts a RethinkDB term.
func (k *K) DB() rt.Term {
	return rt.DB(k.Opts.DBName)
}

// Run is an alias for rethinkdb.Term's Run(), with the Session and
// RunOpts filled in automatically from parameters passed in. Make sure to
// close the cursor!
func (k *K) Run(t rt.Term) (*rt.Cursor, error) {
	return t.Run(k.Session, k.Opts.DBRunOpts)
}

// RunWrite is an alias for rethinkdb.Term's RunWrite(), with the Session and
// RunOpts filled in automatically from parameters passed in, EXCEPT if the
// response completes BUT returns a WriteResponse with the number of errors
// greater than zero, it does NOT return that error as the second return value.
func (k *K) RunWrite(t rt.Term) (rt.WriteResponse, error) {
	var response rt.WriteResponse

	res, err := t.Run(k.Session, k.Opts.DBRunOpts)
	if err != nil {
		return response, err
	}

	if err = res.One(&response); err != nil {
		return response, err
	}

	if err = res.Close(); err != nil {
		return response, err
	}

	return response, nil
}

// RunWriteErr is RunWrite with the added interpretation of WriteResponse's
// errors, which are returned as a WRError if found.
func (k *K) RunWriteErr(t rt.Term) (rt.WriteResponse, error) {
	wr, err := k.RunWrite(t)
	if err != nil {
		return wr, err
	}

	if wr.Errors > 0 {
		return wr, &WRError{wr.FirstError}
	}

	return wr, nil
}
