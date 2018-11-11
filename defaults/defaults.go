// DBDeployer - The MySQL Sandbox
// Copyright © 2006-2018 Giuseppe Maxia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package defaults

import (
	"encoding/json"
	"fmt"
	"github.com/datacharmer/dbdeployer/common"
	"os"
	"path"
	"time"
)

type DbdeployerDefaults struct {
	Version           string `json:"version"`
	SandboxHome       string `json:"sandbox-home"`
	SandboxBinary     string `json:"sandbox-binary"`
	UseSandboxCatalog bool   `json:"use-sandbox-catalog"`
	LogSBOperations   bool   `json:"log-sb-operations"`
	LogDirectory      string `json:"log-directory"`

	//UseConcurrency    			   bool   `json:"use-concurrency"`
	MasterSlaveBasePort           int `json:"master-slave-base-port"`
	GroupReplicationBasePort      int `json:"group-replication-base-port"`
	GroupReplicationSpBasePort    int `json:"group-replication-sp-base-port"`
	FanInReplicationBasePort      int `json:"fan-in-replication-base-port"`
	AllMastersReplicationBasePort int `json:"all-masters-replication-base-port"`
	MultipleBasePort              int `json:"multiple-base-port"`
	// GaleraBasePort                 int    `json:"galera-base-port"`
	// PXCBasePort                    int    `json:"pxc-base-port"`
	// NdbBasePort                    int    `json:"ndb-base-port"`
	GroupPortDelta    int    `json:"group-port-delta"`
	MysqlXPortDelta   int    `json:"mysqlx-port-delta"`
	MasterName        string `json:"master-name"`
	MasterAbbr        string `json:"master-abbr"`
	NodePrefix        string `json:"node-prefix"`
	SlavePrefix       string `json:"slave-prefix"`
	SlaveAbbr         string `json:"slave-abbr"`
	SandboxPrefix     string `json:"sandbox-prefix"`
	MasterSlavePrefix string `json:"master-slave-prefix"`
	GroupPrefix       string `json:"group-prefix"`
	GroupSpPrefix     string `json:"group-sp-prefix"`
	MultiplePrefix    string `json:"multiple-prefix"`
	FanInPrefix       string `json:"fan-in-prefix"`
	AllMastersPrefix  string `json:"all-masters-prefix"`
	ReservedPorts     []int  `json:"reserved-ports"`
	// GaleraPrefix                   string `json:"galera-prefix"`
	// PxcPrefix                      string `json:"pxc-prefix"`
	// NdbPrefix                      string `json:"ndb-prefix"`
	Timestamp string `json:"timestamp"`
}

const (
	minPortValue int = 11000
	maxPortValue int = 30000
)

var (
	homeDir                 string = os.Getenv("HOME")
	ConfigurationDir        string = path.Join(homeDir, ".dbdeployer")
	ConfigurationFile       string = path.Join(ConfigurationDir, "config.json")
	CustomConfigurationFile string = ""
	SandboxRegistry         string = path.Join(ConfigurationDir, "sandboxes.json")
	SandboxRegistryLock     string = path.Join(ConfigurationDir, "sandboxes.lock")
	LogSBOperations         bool   = common.IsEnvSet("DBDEPLOYER_LOGGING")

	// This variable is changed to true when the "cmd" package is activated,
	// meaning that we're using the command line interface of dbdeployer.
	// It is used to make decisions whether to write messages to the screen
	// when calling sandbox creation functions from other apps.
	UsingDbDeployer bool = false

	factoryDefaults = DbdeployerDefaults{
		Version:       common.CompatibleVersion,
		SandboxHome:   path.Join(homeDir, "sandboxes"),
		SandboxBinary: path.Join(homeDir, "opt", "mysql"),

		UseSandboxCatalog: true,
		LogSBOperations:   false,
		LogDirectory:      path.Join(homeDir, "sandboxes", "logs"),
		//UseConcurrency :			   true,
		MasterSlaveBasePort:           11000,
		GroupReplicationBasePort:      12000,
		GroupReplicationSpBasePort:    13000,
		FanInReplicationBasePort:      14000,
		AllMastersReplicationBasePort: 15000,
		MultipleBasePort:              16000,
		// GaleraBasePort:                17000,
		// PxcBasePort:                   18000,
		// NdbBasePort:                   19000,
		GroupPortDelta:    125,
		MysqlXPortDelta:   10000,
		MasterName:        "master",
		MasterAbbr:        "m",
		NodePrefix:        "node",
		SlavePrefix:       "slave",
		SlaveAbbr:         "s",
		SandboxPrefix:     "msb_",
		MasterSlavePrefix: "rsandbox_",
		GroupPrefix:       "group_msb_",
		GroupSpPrefix:     "group_sp_msb_",
		MultiplePrefix:    "multi_msb_",
		FanInPrefix:       "fan_in_msb_",
		AllMastersPrefix:  "all_masters_msb_",
		ReservedPorts:     []int{1186, 3306, 33060},
		// GaleraPrefix:                  "galera_msb_",
		// NdbPrefix:                     "ndb_msb_",
		// PxcPrefix:                     "pxc_msb_",
		Timestamp: time.Now().Format(time.UnixDate),
	}
	currentDefaults DbdeployerDefaults
)

func Defaults() DbdeployerDefaults {
	if currentDefaults.Version == "" {
		if common.FileExists(ConfigurationFile) {
			currentDefaults = ReadDefaultsFile(ConfigurationFile)
		} else {
			currentDefaults = factoryDefaults
		}
	}
	if currentDefaults.LogSBOperations {
		LogSBOperations = true
	}
	return currentDefaults
}

func ShowDefaults(defaults DbdeployerDefaults) {
	defaults = replaceLiteralEnvValues(defaults)
	if common.FileExists(ConfigurationFile) {
		fmt.Printf("# Configuration file: %s\n", ConfigurationFile)
	} else {
		fmt.Println("# Internal values:")
	}
	b, err := json.MarshalIndent(defaults, " ", "\t")
	common.ErrCheckExitf(err, 1, ErrEncodingDefaults, err)
	fmt.Printf("%s\n", b)
}

func WriteDefaultsFile(filename string, defaults DbdeployerDefaults) {
	defaults = replaceLiteralEnvValues(defaults)
	defaultsDir := common.DirName(filename)
	if !common.DirExists(defaultsDir) {
		common.Mkdir(defaultsDir)
	}
	b, err := json.MarshalIndent(defaults, " ", "\t")
	common.ErrCheckExitf(err, 1, ErrEncodingDefaults, err)
	jsonString := fmt.Sprintf("%s", b)
	common.WriteString(jsonString, filename)
}

func expandEnvironmentVariables(defaults DbdeployerDefaults) DbdeployerDefaults {
	defaults.SandboxHome = common.ReplaceEnvVar(defaults.SandboxHome, "HOME")
	defaults.SandboxHome = common.ReplaceEnvVar(defaults.SandboxHome, "PWD")
	defaults.SandboxBinary = common.ReplaceEnvVar(defaults.SandboxBinary, "HOME")
	defaults.SandboxBinary = common.ReplaceEnvVar(defaults.SandboxBinary, "PWD")
	return defaults
}

func replaceLiteralEnvValues(defaults DbdeployerDefaults) DbdeployerDefaults {
	defaults.SandboxHome = common.ReplaceLiteralEnvVar(defaults.SandboxHome, "HOME")
	defaults.SandboxHome = common.ReplaceLiteralEnvVar(defaults.SandboxHome, "PWD")
	defaults.SandboxBinary = common.ReplaceLiteralEnvVar(defaults.SandboxBinary, "HOME")
	defaults.SandboxBinary = common.ReplaceLiteralEnvVar(defaults.SandboxBinary, "PWD")
	return defaults
}

func ReadDefaultsFile(filename string) (defaults DbdeployerDefaults) {
	defaultsBlob := common.SlurpAsBytes(filename)

	err := json.Unmarshal(defaultsBlob, &defaults)
	common.ErrCheckExitf(err, 1, ErrEncodingDefaults, err)
	defaults = expandEnvironmentVariables(defaults)
	return
}

func checkInt(name string, val, min, max int) bool {
	if val >= min && val <= max {
		return true
	}
	fmt.Printf("Value %s (%d) must be between %d and %d\n", name, val, min, max)
	return false
}

func ValidateDefaults(nd DbdeployerDefaults) bool {
	var allInts bool
	allInts = checkInt("master-slave-base-port", nd.MasterSlaveBasePort, minPortValue, maxPortValue) &&
		checkInt("group-replication-base-port", nd.GroupReplicationBasePort, minPortValue, maxPortValue) &&
		checkInt("group-replication-sp-base-port", nd.GroupReplicationSpBasePort, minPortValue, maxPortValue) &&
		checkInt("multiple-base-port", nd.MultipleBasePort, minPortValue, maxPortValue) &&
		checkInt("fan-in-base-port", nd.FanInReplicationBasePort, minPortValue, maxPortValue) &&
		checkInt("all-masters-base-port", nd.AllMastersReplicationBasePort, minPortValue, maxPortValue) &&
		// check_int("galera-base-port", nd.GaleraBasePort, min_port_value, max_port_value) &&
		// check_int("pxc-base-port", nd.PxcBasePort, min_port_value, max_port_value) &&
		// check_int("ndb-base-port", nd.NdbBasePort, min_port_value, max_port_value) &&
		checkInt("group-port-delta", nd.GroupPortDelta, 101, 299)
	checkInt("mysqlx-port-delta", nd.MysqlXPortDelta, 2000, 15000)
	if !allInts {
		return false
	}
	var noConflicts bool
	noConflicts = nd.MultipleBasePort != nd.GroupReplicationSpBasePort &&
		nd.MultipleBasePort != nd.GroupReplicationBasePort &&
		nd.MultipleBasePort != nd.MasterSlaveBasePort &&
		nd.MultipleBasePort != nd.FanInReplicationBasePort &&
		nd.MultipleBasePort != nd.AllMastersReplicationBasePort &&
		// nd.MultipleBasePort != nd.NdbBasePort &&
		// nd.MultipleBasePort != nd.GaleraBasePort &&
		// nd.MultipleBasePort != nd.PxcBasePort &&
		nd.MultiplePrefix != nd.GroupSpPrefix &&
		nd.MultiplePrefix != nd.GroupPrefix &&
		nd.MultiplePrefix != nd.MasterSlavePrefix &&
		nd.MultiplePrefix != nd.SandboxPrefix &&
		nd.MultiplePrefix != nd.FanInPrefix &&
		nd.MultiplePrefix != nd.AllMastersPrefix &&
		nd.MasterAbbr != nd.SlaveAbbr &&
		// nd.MultiplePrefix != nd.NdbPrefix &&
		// nd.MultiplePrefix != nd.GaleraPrefix &&
		// nd.MultiplePrefix != nd.PxcPrefix &&
		nd.SandboxHome != nd.SandboxBinary
	if !noConflicts {
		fmt.Printf("Conflicts found in defaults values:\n")
		ShowDefaults(nd)
		return false
	}
	allStrings := nd.SandboxPrefix != "" &&
		nd.MasterSlavePrefix != "" &&
		nd.MasterName != "" &&
		nd.MasterAbbr != "" &&
		nd.NodePrefix != "" &&
		nd.SlavePrefix != "" &&
		nd.SlaveAbbr != "" &&
		nd.GroupPrefix != "" &&
		nd.GroupSpPrefix != "" &&
		nd.MultiplePrefix != "" &&
		nd.SandboxHome != "" &&
		nd.SandboxBinary != ""
	if !allStrings {
		fmt.Printf("One or more empty values found in defaults\n")
		ShowDefaults(nd)
		return false
	}
	versionList := common.VersionToList(common.CompatibleVersion)
	if !common.GreaterOrEqualVersion(nd.Version, versionList) {
		fmt.Printf("Provided defaults are for version %s. Current version is %s\n", nd.Version, common.CompatibleVersion)
		return false
	}
	return true
}

func RemoveDefaultsFile() {
	if common.FileExists(ConfigurationFile) {
		err := os.Remove(ConfigurationFile)
		common.ErrCheckExitf(err, 1, "%s", err)
		fmt.Printf("#File %s removed\n", ConfigurationFile)
	} else {
		common.Exitf(1, "configuration file %s not found", ConfigurationFile)
	}
}

func strToSlice(label, s string) []int {
	intList, err := common.StringToIntSlice(s)
	if err != nil {
		common.Exitf(1, "bad input for %s: %s (%s) ", label, s, err)
	}
	return intList
}

func UpdateDefaults(label, value string, storeDefaults bool) {
	newDefaults := Defaults()
	switch label {
	case "version":
		newDefaults.Version = value
	case "sandbox-home":
		newDefaults.SandboxHome = value
	case "sandbox-binary":
		newDefaults.SandboxBinary = value
	case "use-sandbox-catalog":
		newDefaults.UseSandboxCatalog = common.TextToBool(value)
	case "log-sb-operations":
		newDefaults.LogSBOperations = common.TextToBool(value)
	case "log-directory":
		newDefaults.LogDirectory = value
	//case "use-concurrency":
	//	new_defaults.UseConcurrency = common.TextToBool(value)
	case "master-slave-base-port":
		newDefaults.MasterSlaveBasePort = common.Atoi(value)
	case "group-replication-base-port":
		newDefaults.GroupReplicationBasePort = common.Atoi(value)
	case "group-replication-sp-base-port":
		newDefaults.GroupReplicationSpBasePort = common.Atoi(value)
	case "multiple-base-port":
		newDefaults.MultipleBasePort = common.Atoi(value)
	case "fan-in-base-port":
		newDefaults.FanInReplicationBasePort = common.Atoi(value)
	case "all-masters-base-port":
		newDefaults.AllMastersReplicationBasePort = common.Atoi(value)
	// case "ndb-base-port":
	//	 new_defaults.NdbBasePort = common.Atoi(value)
	// case "galera-base-port":
	//	 new_defaults.GaleraBasePort = common.Atoi(value)
	// case "pxc-base-port":
	//	 new_defaults.PxcBasePort = common.Atoi(value)
	case "group-port-delta":
		newDefaults.GroupPortDelta = common.Atoi(value)
	case "mysqlx-port-delta":
		newDefaults.MysqlXPortDelta = common.Atoi(value)
	case "master-name":
		newDefaults.MasterName = value
	case "master-abbr":
		newDefaults.MasterAbbr = value
	case "node-prefix":
		newDefaults.NodePrefix = value
	case "slave-prefix":
		newDefaults.SlavePrefix = value
	case "slave-abbr":
		newDefaults.SlaveAbbr = value
	case "sandbox-prefix":
		newDefaults.SandboxPrefix = value
	case "master-slave-prefix":
		newDefaults.MasterSlavePrefix = value
	case "group-prefix":
		newDefaults.GroupPrefix = value
	case "group-sp-prefix":
		newDefaults.GroupSpPrefix = value
	case "multiple-prefix":
		newDefaults.MultiplePrefix = value
	case "fan-in-prefix":
		newDefaults.FanInPrefix = value
	case "all-masters-prefix":
		newDefaults.AllMastersPrefix = value
	case "reserved-ports":
		newDefaults.ReservedPorts = strToSlice("reserved-ports", value)
	// case "galera-prefix":
	// 	new_defaults.GaleraPrefix = value
	// case "pxc-prefix":
	// 	new_defaults.PxcPrefix = value
	// case "ndb-prefix":
	// 	new_defaults.NdbPrefix = value
	default:
		common.Exitf(1, "unrecognized label %s", label)
	}
	if ValidateDefaults(newDefaults) {
		currentDefaults = newDefaults
		if storeDefaults {
			WriteDefaultsFile(ConfigurationFile, Defaults())
			fmt.Printf("# Updated %s -> \"%s\"\n", label, value)
		}
	} else {
		common.Exitf(1, "invalid defaults data %s : %s", label, value)
	}
}

func LoadConfiguration() {
	if !common.FileExists(ConfigurationFile) {
		// WriteDefaultsFile(ConfigurationFile, Defaults())
		return
	}
	newDefaults := ReadDefaultsFile(ConfigurationFile)
	if ValidateDefaults(newDefaults) {
		currentDefaults = newDefaults
	} else {
		fmt.Println(common.StarLine)
		fmt.Printf("Defaults file %s not validated.\n", ConfigurationFile)
		fmt.Println("Loading internal defaults")
		fmt.Println(common.StarLine)
		fmt.Println("")
		time.Sleep(1000 * time.Millisecond)
	}
}
