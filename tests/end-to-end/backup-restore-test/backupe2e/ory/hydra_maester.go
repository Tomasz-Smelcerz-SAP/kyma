package ory

import (
	"log"
	"time"

	"github.com/avast/retry-go"
	"github.com/Tomasz-Smelcerz-SAP/kyma/tests/integration/api-gateway/gateway-tests/pkg/client"
	"github.com/Tomasz-Smelcerz-SAP/kyma/tests/integration/api-gateway/gateway-tests/pkg/manifestprocessor"
	"github.com/Tomasz-Smelcerz-SAP/kyma/tests/integration/api-gateway/gateway-tests/pkg/resource"
	. "github.com/smartystreets/goconvey/convey"
)

const hydraClientFile = "../backupe2e/ory/manifests/hydra-client.yaml"
const manifestsDirectory = "manifests"
const resourceSeparator = "---"

var resourceManager *resource.Manager
var batch *resource.Batch

type HydraClientTest struct {
}

type hydraClientScenario struct {
	namespace string
}

type scenarioStep func() error

type Config struct {
	HydraAddr  string `envconfig:"TEST_HYDRA_ADDRESS"`
	User       string `envconfig:"TEST_USER_EMAIL"`
	Pwd        string `envconfig:"TEST_USER_PASSWORD"`
	ReqTimeout uint   `envconfig:"TEST_REQUEST_TIMEOUT,default=100"`
	ReqDelay   uint   `envconfig:"TEST_REQUEST_DELAY,default=5"`
	Domain     string `envconfig:"DOMAIN"`
}

func NewHydraClientTest() (*HydraClientTest, error) {

	conf := Config{
		ReqDelay:   4,
		ReqTimeout: 12,
	}

	commonRetryOpts := []retry.Option{
		retry.Delay(time.Duration(conf.ReqDelay) * time.Second),
		retry.Attempts(conf.ReqTimeout / conf.ReqDelay),
		retry.DelayType(retry.FixedDelay),
	}

	resourceManager = &resource.Manager{RetryOptions: commonRetryOpts}
	batch = &resource.Batch{
		resourceManager,
	}

	return &HydraClientTest{}, nil
}

func (hct *HydraClientTest) CreateResources(namespace string) {
	hct.newScenario(namespace).createResources()
}

func (hct *HydraClientTest) TestResources(namespace string) {
	hct.newScenario(namespace).testResources()
}

func (hct *HydraClientTest) newScenario(namespace string) *hydraClientScenario {
	return &hydraClientScenario{namespace}
}

func (hcs *hydraClientScenario) createResources() {
	hcs.run(
		hcs.createClient,
		//hcs.fetchToken,
		//hcs.verifyToken,
	)
}

func (hcs *hydraClientScenario) testResources() {
	hcs.run(
		hcs.readClient,
		hcs.fetchToken,
		hcs.verifyToken,
	)
}

func (hcs *hydraClientScenario) run(steps ...scenarioStep) {
	for _, fn := range steps {
		err := fn()
		if err != nil {
			log.Println(err)
		}
		So(err, ShouldBeNil)
	}
}

func (hcs *hydraClientScenario) createClient() error {
	//manifestprocessor.ParseFromFileWithTemplate(testingAppFile, manifestsDirectory, resourceSeparator, struct{ TestID string }{TestID: testID})
	hydraClientResource, err := manifestprocessor.ParseFromFileWithTemplate(hydraClientFile, manifestsDirectory, resourceSeparator, struct{ TestNamespace string }{TestNamespace: hcs.namespace})

	if err != nil {
		panic(err)
	}
	log.Println(hydraClientResource)

	k8sClient := client.GetDynamicClient()
	batch.CreateResources(k8sClient, hydraClientResource...)

	return nil
}

func (hcs *hydraClientScenario) readClient() error {
	return nil
}

func (hcs *hydraClientScenario) fetchToken() error {
	return nil
}

func (hcs *hydraClientScenario) verifyToken() error {
	return nil
}
