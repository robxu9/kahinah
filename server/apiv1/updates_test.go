package apiv1

import (
	"net/http"
	"testing"
)

func TestUpdates(t *testing.T) {
	setup := setupAPIv1(t)
	defer setup.Close()

	setup.MakeUpdate()

	wantedResponse := `{"links":{"head":"/updates","tail":"/updates","total":"1"},"updates":[{"Id":1,"AdvisoryId":0,"For":"robxu9/2014/main","Name":"test-1.0.0/amd64","Submitter":"test@example.com","Type":1,"Content":{"Id":1,"From":"120","To":"126","Url":"http://example.com/","BuiltAt":"2014-04-30T20:00:00-04:00","Packages":[{"Name":"test","Epoch":0,"Version":"1.0.0","Release":"1.robxu9","Arch":"amd64","Type":"src","Url":"http://example.com/test.tar.xz"},{"Name":"test","Epoch":0,"Version":"1.0.0","Release":"1.robxu9","Arch":"amd64","Type":"binary","Url":"http://example.com/test.pkg"}],"Changes":[{"ChangeAt":"2014-04-30T20:00:00-04:00","For":"1.0.0-1.robxu9","By":"test@example.com","Details":"did some stuff"}]},"CreatedAt":"*"}]}`

	setup.NetTest("GET", "/updates", "", false, http.StatusOK, wantedResponse, MatchGlob)
}
