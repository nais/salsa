package jvm

import (
	"reflect"
	"testing"
)

func TestMavenDeps(t *testing.T) {
	got, _ := MavenCompileAndRuntimeTimeDeps(mvnDependencyListOutput)
	want := map[string]string{
		"org.apache.logging.log4j:log4j-core": "2.14.1",
		"org.apache.logging.log4j:log4j-api":  "2.14.2",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}

}

const mvnDependencyListOutput = `[INFO] Scanning for projects...
[INFO]
[INFO] -----------------------------< jks:stuff >------------------------------
[INFO] Building stuff 1.0-SNAPSHOT
[INFO] --------------------------------[ jar ]---------------------------------
[INFO]
[INFO] --- maven-dependency-plugin:2.8:list (default-cli) @ stuff ---
[WARNING] The artifact xml-apis:xml-apis:jar:2.0.2 has been relocated to xml-apis:xml-apis:jar:1.0.b2
[INFO]
[INFO] The following files have been resolved:
[INFO]    junit:junit:jar:3.8.1:test
[INFO]    org.apache.logging.log4j:log4j-core:jar:2.14.1:compile
[INFO]    org.apache.logging.log4j:log4j-api:jar:2.14.2:compile
[INFO]
[INFO] ------------------------------------------------------------------------
[INFO] BUILD SUCCESS
[INFO] ------------------------------------------------------------------------
[INFO] Total time:  0.864 s
[INFO] Finished at: 2021-10-20T19:31:08+02:00
[INFO] ------------------------------------------------------------------------`
