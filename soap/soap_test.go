package soap

import (
	"context"
	"encoding/xml"
	"github.com/evorts/kevlars/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/expect"
	"github.com/vitorsalgado/mocha/v3/reply"
	"net/http"
	"strings"
	"testing"
)

type TestSuite struct {
	suite.Suite
	m    Manager
	mock *mocha.Mocha
}

func (ts *TestSuite) SetupTest() {
	ts.m = New(
		WithBasicAuth("aladin_demo", "Aladin2023!@"),
		WithTransport(WithTelemetry(
			telemetry.NewNoop().MustInit(),
		)),
	)
	ts.mock = mocha.New(ts.T())
	ts.mock.Start()
}

func (ts *TestSuite) TestSoapRaw() {
	rawReq := []byte(strings.TrimSpace(`
<?xml version="1.0" encoding="UTF-8"?>
<env:Envelope xmlns:env="http://www.w3.org/2001/12/soap-envelope">
		<env:Header/>
		<env:Body>
				<getPetByIdRequest xmlns="urn:com:example:petstore">
						<id>3</id>
				</getPetByIdRequest>
		</env:Body>
</env:Envelope>
`))
	ctx := context.Background()
	var rs struct {
		XMLName xml.Name `xml:"Envelope"`
		Text    string   `xml:",chardata"`
		Env     string   `xml:"env,attr"`
		Header  string   `xml:"Header"`
		Body    struct {
			Text               string `xml:",chardata"`
			GetPetByIdResponse struct {
				Text  string `xml:",chardata"`
				Xmlns string `xml:"xmlns,attr"`
				ID    string `xml:"id"`
				Name  string `xml:"name"`
			} `xml:"getPetByIdResponse"`
		} `xml:"Body"`
	}
	scoped := ts.mock.AddMocks(
		mocha.Post(expect.URLPath("/soap/")).
			Header("Content-Type", expect.ToEqual("application/xml").Or(expect.ToEqual("text/xml"))).
			Reply(reply.OK().Body([]byte(strings.TrimSpace(`
<?xml version="1.0" encoding="UTF-8"?>
<env:Envelope
	xmlns:env="http://www.w3.org/2001/12/soap-envelope">
	<env:Header/>
	<env:Body>
		<getPetByIdResponse
			xmlns="urn:com:example:petstore">
			<id>3</id>
			<name>Pet Name</name>
		</getPetByIdResponse>
	</env:Body>
</env:Envelope>`)))),
	)
	httpCode, err := ts.m.PostRaw(
		ctx, "",
		ts.mock.URL()+"/soap/",
		rawReq, &rs,
	)
	require.NoError(ts.T(), err)
	assert.Nil(ts.T(), err)
	assert.True(ts.T(), scoped.Called())
	assert.Equal(ts.T(), http.StatusOK, httpCode)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
