package main 

import (
	"fmt"
	"os"
	"time"
	"flag"
	"encoding/json"
  "database/sql"
	_ "github.com/lib/pq"
	"archive/tar"
	"path/filepath"
	"strings"
	"regexp"
	"os/exec"
	"encoding/hex"
)
var (
	//sudo -u tempest-web psql -U tempest-web -d tempest_production
	
	// PGHOST is the hostname of opsmanager 
	PGHOST = "127.0.0.1" // TODO not used 
	// PGUSER is the opsmanager database username
	PGUSER = "tempest-web"
	// PGDATABASE is the opsmanager database name
	PGDATABASE = "tempest_production"
	// PGPORT is the port for the opsman database
	PGPORT = "5432" // TODO  not used 
	// PGPASS is optional and is the password for the opsmanager database user
	PGPASS = ""
	// OPSDBTYPE deafults to postgress but can be changed via environment variable too
	OPSDBTYPE = "postgres" 
	// DBURL is the connection string used to establish database session
	DBURL = ""
	
	//Collors
	textColor = "\033[94m"
	errorColor = "\033[91m"
	warnColor = "\033[33m"
	termColor = "\033[0m"
	
	gzipCMD = "gzip"
	
	changesLogQuery = `SELECT * FROM installation_changes`
	changeLogDataQuery = `SELECT id, created_at, updated_at, install_id, encode(text, 'hex') from installation_logs`

	customer = flag.String("c", "ANONYMOUS", "Name of customer")
	prodenv = flag.Bool("p", false, "Use this flag if this is a prod environemnt")
	testenv = flag.Bool("n", false, "Use this flag if this is a non-prod environment")
	outputdir = flag.String("o", "./", "defaults to current working director")
)

// ChangeLogChanges database table struct for installation_changes table
type ChangeLogChanges struct {
	ID int `json:"id"`
	Customer string `json:"customer"`
	Identifier string `json:"indentifier"`
	Label string `json:"label"`
	GUID string `json:"guid"`
	ProductVersion string `json:"product_version"`// opsman version change  
	InstallID int `json:"install_id"`
	ChangeType string `json:"change_type"`// addition, update, ..
}

// ChangeLogs database table struct for installation_logs table
type ChangeLogs struct {
	ID int `json:"id"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
	InstallID int `json:"install_id"`
	Log []byte `json:"log"`
}

// Check environment for postgres environment variable and override defaults
func checkEnv(s *string, key string) {
	c := os.Getenv(key)
	if c != "" {
		*s = c
	}
}

// Marshal the given struct into json format
func marshalStruct(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err !=nil {
		panic(fmt.Sprintf("Could not marshal struct:%s", err))
	}
	return b
}

// read in data from the database
func collect(outdir, outfile string) error {
	sess, err := sql.Open(OPSDBTYPE, DBURL)
	if err != nil {
		return fmt.Errorf("Failed to connect to database %s: %s", DBURL, err)
	}	
	defer sess.Close()
	
	changes := make([]ChangeLogChanges,0)
	change := ChangeLogChanges{}
	log := ChangeLogs{}
	tarfh, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	tw := tar.NewWriter(tarfh)
	defer tarfh.Close()
	defer tw.Close()
	
	rows, err := sess.Query(changesLogQuery)
	if err != nil {
		return fmt.Errorf("Query Failed: \"%s\": %s", changesLogQuery, err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&change.ID,
		&change.Identifier,
		&change.Label,
		&change.GUID,
		&change.ProductVersion,
		&change.InstallID, 
		&change.ChangeType)
		if err != nil {
			return fmt.Errorf("Row Scan Failed for changes: %s", err)
		}
		
		//TODO need to write this data out to tar file
		changes = append(changes, change)
	}
	
	changesData := marshalStruct(changes)
	err = tw.WriteHeader(&tar.Header{
			Name: "changes.txt",
			Mode: 0666,
			Size: int64(len(changesData)),
		})
	if err != nil {
		return fmt.Errorf("Failed write tar header for changes:%d: %s", int64(len(changesData)), err)
	}
	
	_, err = tw.Write(changesData)
	if err != nil {
		return fmt.Errorf("Failed to write changes file to archive: %s", err)
	}
	
	// Get the log data an create a file in the archive for each row returned
	logRows, err := sess.Query(changeLogDataQuery)
	if err != nil {
		return fmt.Errorf("Query Failed: \"%s\": %s", changeLogDataQuery, err)
	}
	defer logRows.Close()
	for logRows.Next() {
		err = logRows.Scan(&log.ID,
			&log.CreatedAT,
			&log.UpdatedAT,
			&log.InstallID,
			&log.Log)
		if err != nil {
			return fmt.Errorf("Row Scan Failed for logs: %s", err)
		}
		//b := marshalStruct(log)
		b, err := hex.DecodeString(fmt.Sprintf("%s", log.Log))
		if err != nil {
			return fmt.Errorf("Unable to decode hex log into []bytes: %s", err)
		}
		err = tw.WriteHeader(&tar.Header{
				Name: fmt.Sprintf("%d_%d_%d_%d_changelog.txt", log.ID, log.InstallID, log.CreatedAT.Unix(), log.UpdatedAT.Unix()), 
				Mode: 0666,
				Size: int64(len(b)),
			})
		if err != nil {
			return fmt.Errorf("Failed write tar header for log file:%d: %s", int64(len(b)), err)
		}
		_, err = tw.Write(b)
		if err != nil {
			return fmt.Errorf("Failed to write changes file to archive: %s", err)
		}
	}
	return nil
}

// remove spaces and special characeters and change from uppercase to lowercase
func cleanCustomerName(n *string) {
	re := regexp.MustCompile(`[a-z]|[0-9]`)
	*n = strings.ToLower(*n)
	*n = strings.Join(re.FindAllString(*n, -1), "")
}

func resolveCommands() {
	var err error
	gzipCMD, err = exec.LookPath(gzipCMD)
	if err != nil {
		panic(fmt.Sprintf("Problem with gzip command dependencies: %s", err))
	}
}

func logError(s string){
	fmt.Printf("~# %s%s%s\n", errorColor, s, termColor)
}
func logWarn(s string){
	fmt.Printf("~# %s%s%s\n", warnColor, s, termColor)
}
func logInfo(s string){
	fmt.Printf("~# %s%s%s\n", textColor, s, termColor)
}

func main(){
	checkEnv(&PGHOST, "PGHOST")
	checkEnv(&PGUSER, "PGUSER")
	checkEnv(&PGDATABASE, "PGDATABASE")
	checkEnv(&PGPORT, "PGPORT")
	checkEnv(&PGPASS, "PGPASS")
	checkEnv(&OPSDBTYPE, "OPSDBTYPE")
	flag.Parse()
	
	envType := ""
	if ! *prodenv && ! *testenv {
		logError("Please specify -p if this is a prod environment or -n if this is a non-prod environment")
		os.Exit(2)
	}
	if *prodenv {
		envType = "prod"
	} else if *testenv {
		envType = "non-prod"
	}
	
	if *customer == "ANONYMOUS" {
		logWarn("Warning customer name will default to ANONYMOUS unless you specify the name with -c option")
	}
	
	//DBURL = fmt.Sprintf("%s://%s@%s:%s/%s", OPSDBTYPE, PGUSER, PGHOST, PGPORT, PGDATABASE)
	DBURL = fmt.Sprintf("%s:///%s?host=/var/run/postgresql/", OPSDBTYPE, PGDATABASE)
	
	cleanCustomerName(customer)
	logInfo(fmt.Sprintf("Scrubbed customer name and will use %s", *customer))
	outfile := filepath.Join(*outputdir, fmt.Sprintf("%s_%s_opsman-changelogs.tar", *customer, envType))
	logInfo(fmt.Sprintf("Creating tar archive %s", outfile))
	err := collect(*outputdir, outfile)
	if err != nil {
		logError(err.Error())
		os.Exit(1)
	}
	
	logInfo("Compressing archive")
	out, err := exec.Command(gzipCMD, outfile).CombinedOutput()
	if err != nil {
		logError(fmt.Sprintf("Failed to Compress Archive %s: %s", out, err))
		os.Exit(3)
	}
	logInfo(fmt.Sprintf("Successfully crated archive %s.gz", outfile))
}
