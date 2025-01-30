package coordinator_client

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	os.Exit(m.Run())
}

func TestRedisCoordinatorClientCreateTask(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	taskParams := map[string]string{
		"url": "https://ethanhosier.com",
	}

	task, err := NewTask("2b665be2-80b7-40d4-9117-a6e9794afe97", "asdasdasd", taskParams)
	if err != nil {
		t.Fatalf("Failed to    create ta sk:   %v", err)
	}

	err = client.CreateTask(context.Background(), CoordinatorClientTaskTopicUrls, task)
	if err != nil {
		t.Fatalf("Failed to create task      : %v", err)
	}
}

func TestRedisCoordinatorClientCreate100Tasks(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	for i := 0; i < 100; i++ {
		params := map[string]string{
			"url": "https://ethanhosier.com?test=" + strconv.Itoa(i),
		}

		task, err := NewTask(uuid.New().String(), "asdasdasd", params)
		if err != nil {
			t.Fatalf("Failed to  create  task: %v", err)
		}

		err = client.CreateTask(context.Background(), CoordinatorClientTaskTopicUrls, task)
		if err != nil {
			t.Fatalf("Failed to crea te task: %v", err)
		}
	}
}

func TestRedisCoordinatorClientCreate100UrlTasks(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	links := `carrtool.com
dacres.co.uk
chalets-village.com
allenslanding.ca
ghanja.be
luzcasal.es
domaine-du-jasson.com
asp.es
smarterfinancialplanning.ca
cal-inc.com
sintnicolaashoeve.nl
vitalingua.de
sports-hoop.com
abs.edu
lennyssurfshop.com
icucciolidicasapapa.it
nicolasbaleydier.com
helloeye.co.kr
agilepm.com
ritualbedarf.de
iccsl.es
borgodeiferraresi.com
tenerife-online.ro
aikidomadrid.cl
schwimmshop.de
kamo-photoart.de
rogerssportinggoods.com
cgscgs.com
tlmassociates.com
textkritik.de
jimo.gov.cn
aversi.ge
cyrusimap.org
mscgis.de
thetaylormadeteam.com
greenriverfestival.com
abrams-california-health-insurance.com
harlowjudoclub.com
jbeelsdesign.com
acrf.org
worldjidokwan.com
appliancedoctor.co.uk
deshihost.com
myudm.ru
beacon.org
gchbuilders.com
norm-flynn.com
provincia.siracusa.it
ruthsager.com.ar
iwate-safari.jp
illusion.ee
aurochs.org
room-service-bonn.de
sign-central.com
poezenweide.nl
labastideduroy.com
withambaptist.org.uk
ericscottmusic.com
tv-wiblingen.de
scorpio.pl
pensiunea-florina.ro
security-doors-windows.com.au
maxhale.nu
nedis.nl
win-rar.ru
missmoss.de
rodella.com
therealgarden.co.uk
grossostheim.de
nienaber-uhren.de
arizonachamberexecs.com
esteticabenessere.it
cpta.org
kgpsoftware.com
automobile-gentz.de
marlboroughrunningclub.org.uk
kaledonya.com
opt-net.com
syedbawkher.com
westendbad.de
cnwrecovery.co.uk
danoneinstitute.org
imkerverein.at
nayada.by
sv-ohmbach.de
bwbcpa.com
925.co.uk
onecommunityfcu.org
mircscripts.fr
calangianus.ot.it
shloimedachs.com
plattweb.co.uk
ryukyu-kyusho.com
fdu.org.ua
virus-protect.org
vendingtimes.com
goergen-rechtsanwalt.de
accounting-by-post.com
durian.in
ideabuilders.com`

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)
	urls := strings.Split(links, "\n")

	for _, url := range urls {
		params := map[string]string{
			"url": url,
		}

		task, err := NewTask(uuid.New().String(), "asdasdsasd", params)
		if err != nil {
			t.Fatalf("Failed t o    create task: %v", err)
		}

		err = client.CreateTask(context.Background(), CoordinatorClientTaskTopicUrls, task)
		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}
	}
}

func TestRedisCoordinatorClientGetTask(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}
	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	task, err := client.GetTask(context.Background(), 5*time.Second, CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	t.Logf("Task: %+v", task)
}

func TestRedisCoordinatorClientGetTaskAndSetProcessing(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	task, err := client.GetTaskAndSetProcessing(context.Background(), 5*time.Second, CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	t.Logf("Task: %+v", task)
}

func TestRedisCoordinatorClientSetProcessed(t *testing.T) {
	task, err := NewTask("37407602-a309-4afd-8b77-efa91d808bf3", "asdasdasd", "test")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD ")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	client.SetProcessed(context.Background(), CoordinatorClientTaskTopicUrls, task)
}

func TestRedisCoordinatorClientCreateGetTask(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	type TestStruct struct {
		Number int    `json:"number"`
		Name   string `json:"name"`
	}

	params := TestStruct{Number: 1, Name: "a name"}

	task, err := NewTask("37407602-a309-4afd-8b77-efa91d808bf3", "asdasdasd", params)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	err = client.CreateTask(context.Background(), CoordinatorClientTaskTopicUrls, task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	task, err = client.GetTask(context.Background(), 5*time.Second, CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	parsedParams, err := CastParams[TestStruct](task.Params)
	if err != nil {
		t.Fatalf("Failed to parse params: %v", err)
	}

	assert.Equal(t, parsedParams.Number, params.Number)
	assert.Equal(t, parsedParams.Name, params.Name)
}

func TestRedisCoordinatorClientStoreError(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	params := map[string]string{
		"url": "https://ethanhosier.com",
	}

	task, err := NewTask("37407602-a309-4afd-8b77-efa91d808bf3", "asdasdasd", params)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	err = client.StoreError(context.Background(), CoordinatorClientTaskTopicUrls, task, fmt.Errorf("an error"))
	if err != nil {
		t.Fatalf("Failed to store error : %v", err)
	}
}
