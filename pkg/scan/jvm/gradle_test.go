package jvm

import (
	"reflect"
	"testing"
)

func TestGradleDeps(t *testing.T) {
	got, _ := GradleDeps(gradleDependenciesOutput)
	want := map[string]string{
		"ch.qos.logback:logback-classic":                          "1.2.9",
		"org.apache.logging.log4j:log4j-core":                     "2.14.1",
		"org.apache.logging.log4j:log4j-api":                      "2.14.2",
		"org.wso2.orbit.org.zapache.commons:commons-collections4": "4.4.wso2v1",
		"com.nimbusds:oauth2-oidc-sdk":                            "9.17",
		"com.github.stephenc.jcip:jcip-annotations":               "1.0-1",
		"junit:junit": "3.8.1",
	}

	if !reflect.DeepEqual(got.Deps, want) {
		t.Errorf("got %q, wanted %q", got.Deps, want)
	}
}

const gradleDependenciesOutput = `testCompileClasspath - Compile classpath for source set 'test'.
+--- ch.qos.logback:logback-classic -> 1.2.9
+--- org.apache.logging.log4j:log4j-core:2.14.1
|    \--- org.apache.logging.log4j:log4j-api:2.14.2
+--- org.wso2.orbit.org.zapache.commons:commons-collections4:4.4.wso2v1
+--- com.nimbusds:oauth2-oidc-sdk:9.17
|    +--- com.github.stephenc.jcip:jcip-annotations:1.0-1
\--- junit:junit:3.8.1`
