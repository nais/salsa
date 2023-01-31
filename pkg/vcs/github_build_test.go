package vcs

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnvironmentGetCurrentFilteredEnvironment(t *testing.T) {
	filteredResult := toEnvData(t, envs)
	assert.Equal(t, 13, len(filteredResult))
	expected := toEnvData(t, expectedEnvs)
	assert.Equal(t, 3, len(expected))
}

func TestParseBuildContext(t *testing.T) {
	expected := make(map[string]string)
	expected["GOVERSION"] = "1.17"
	expected["GOROOT"] = "/opt/hostedtoolcache/go/1.17.6/x64"
	encodedContext := base64.StdEncoding.EncodeToString([]byte(envContext))
	env, err := ParseBuild(&encodedContext)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(env.GetEnvs()))
}

func TestRemoveDuplicateValues(t *testing.T) {
	encodedEnvs := base64.StdEncoding.EncodeToString([]byte(expectedDuplicateEnvs))
	env, err := ParseBuild(&encodedEnvs)
	assert.Equal(t, 5, len(env.GetEnvs()))
	assert.NoError(t, err)
	env.removeDuplicateValues()
	assert.Equal(t, 1, len(env.GetEnvs()))
}

var expectedDuplicateEnvs = `{
        "CLOUDSDK_CORE_PROJECT": "plattformsikkerhet-dev-496e",
        "CLOUDSDK_PROJECT": "plattformsikkerhet-dev-496e",
        "GCLOUD_PROJECT": "plattformsikkerhet-dev-496e",
        "GCP_PROJECT": "plattformsikkerhet-dev-496e",
        "GOOGLE_CLOUD_PROJECT": "plattformsikkerhet-dev-496e"
}`

func TestParseBuildFailEnvironmentalData(t *testing.T) {
	data := "yolo"
	env, err := ParseBuild(&data)
	assert.Nil(t, env)
	assert.EqualError(t, err, "unmarshal environmental context json: invalid character 'Ê' looking for beginning of value")
}

func TestParseBuildNoEnvironmentalData(t *testing.T) {
	env := toEnvData(t, `{}`)
	assert.Equal(t, map[string]string{}, env)
}

func toEnvData(t *testing.T, inputEnvs string) map[string]string {
	encodedEnvs := base64.StdEncoding.EncodeToString([]byte(inputEnvs))
	env, err := ParseBuild(&encodedEnvs)
	assert.NoError(t, err)
	return env.FilterEnvs()
}

var envs = `
{
        "ACTIONS_CACHE_URL": "https://artifactcache.actions.githubusercontent.com/G5IzeHZIHi3eXjyACzx6HDFSkZH4w1o4TXsvvVKrrH5pHDOBRI/",
        "ACTIONS_RUNTIME_TOKEN": "ey..",
        "ACTIONS_RUNTIME_URL": "https://pipelines.actions.githubusercontent.com/G5IzeHZIHi3eXjyACzx6HDFSkZH4w1o4TXsvvVKrrH5pHDOBRI/",
        "CHARSET": "UTF-8",
        "CI": "true",
        "CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE": "/github/workspace/gha-creds-f1d01cf9a81874ff.json",
        "CLOUDSDK_CORE_PROJECT": "plattformsikkerhet-dev-496e",
        "CLOUDSDK_PROJECT": "plattformsikkerhet-dev-496e",
        "GCLOUD_PROJECT": "plattformsikkerhet-dev-496e",
        "GCP_PROJECT": "plattformsikkerhet-dev-496e",
        "GOOGLE_APPLICATION_CREDENTIALS": "/github/workspace/gha-creds-f1d01cf9a81874ff.json",
        "GOOGLE_CLOUD_PROJECT": "plattformsikkerhet-dev-496e",
        "GOOGLE_GHA_CREDS_PATH": "/github/workspace/gha-creds-f1d01cf9a81874ff.json",
        "HOME": "/github/home",
        "HOSTNAME": "85edc7da328e",
        "IMAGE": "ttl.sh/nais/salsa-test:1h",
        "INPUT_GITHUB_CONTEXT": "{\n  \"token\": \"ggg\",\n  \"job\": \"generate-provenance\",\n  \"ref\": \"refs/heads/main\",\n  \"sha\": \"098c4ca5f532e0ffaf9baa19b89bec9905ae6aeb\",\n  \"repository\": \"nais/salsa\",\n  \"repository_owner\": \"nais\",\n  \"repositoryUrl\": \"git://github.com/nais/salsa.git\",\n  \"run_id\": \"1945930372\",\n  \"run_number\": \"60\",\n  \"retention_days\": \"90\",\n  \"run_attempt\": \"1\",\n  \"actor\": \"johnDoe\",\n  \"workflow\": \"NAIS SALSA\",\n  \"head_ref\": \"\",\n  \"base_ref\": \"\",\n  \"event_name\": \"push\",\n  \"event\": {\n    \"after\": \"098c4ca5f532e0ffaf9baa19b89bec9905ae6aeb\",\n    \"base_ref\": null,\n    \"before\": \"4882115251a96ff37321d250a8522d371f5c1720\",\n    \"commits\": [\n      {\n        \"author\": {\n          \"email\": \"john.doe@email.me\",\n          \"name\": \"johnDoe\",\n          \"username\": \"johnDoe\"\n        },\n        \"committer\": {\n          \"email\": \"john.doe@email.me\",\n          \"name\": \"johnDoe\",\n          \"username\": \"johnDoe\"\n        },\n        \"distinct\": true,\n        \"id\": \"098c4ca5f532e0ffaf9baa19b89bec9905ae6aeb\",\n        \"message\": \"add: env to NonRep metadata object\",\n        \"timestamp\": \"2022-03-07T15:13:06+01:00\",\n        \"tree_id\": \"0f3e5162e8131ab2e85190ebd2b2597136b50b9c\",\n        \"url\": \"https://github.com/nais/salsa/commit/098c4ca5f532e0ffaf9baa19b89bec9905ae6aeb\"\n      }\n    ],\n    \"compare\": \"https://github.com/nais/salsa/compare/4882115251a9...098c4ca5f532\",\n    \"created\": false,\n    \"deleted\": false,\n    \"enterprise\": {\n      \"avatar_url\": \"https://avatars.githubusercontent.com/b/371?v=4\",\n      \"created_at\": \"2019-06-26T11:17:54Z\",\n      \"description\": \"\",\n      \"html_url\": \"https://github.com/enterprises/nav\",\n      \"id\": 371,\n      \"name\": \"NAV\",\n      \"node_id\": \"MDEwOkVudGVycHJpc2UzNzE=\",\n      \"slug\": \"nav\",\n      \"updated_at\": \"2021-11-10T10:20:48Z\",\n      \"website_url\": \"https://nav.no\"\n    },\n    \"forced\": false,\n    \"head_commit\": {\n      \"author\": {\n        \"email\": \"john.doe@email.me\",\n        \"name\": \"johnDoe\",\n        \"username\": \"johnDoe\"\n      },\n      \"committer\": {\n        \"email\": \"john.doe@email.me\",\n        \"name\": \"johnDoe\",\n        \"username\": \"johnDoe\"\n      },\n      \"distinct\": true,\n      \"id\": \"098c4ca5f532e0ffaf9baa19b89bec9905ae6aeb\",\n      \"message\": \"add: env to NonRep metadata object\",\n      \"timestamp\": \"2022-03-07T15:13:06+01:00\",\n      \"tree_id\": \"0f3e5162e8131ab2e85190ebd2b2597136b50b9c\",\n      \"url\": \"https://github.com/nais/salsa/commit/098c4ca5f532e0ffaf9baa19b89bec9905ae6aeb\"\n    },\n    \"organization\": {\n      \"avatar_url\": \"https://avatars.githubusercontent.com/u/29488289?v=4\",\n      \"description\": \"NAV Application Infrastructure Service\",\n      \"events_url\": \"https://api.github.com/orgs/nais/events\",\n      \"hooks_url\": \"https://api.github.com/orgs/nais/hooks\",\n      \"id\": 29488289,\n      \"issues_url\": \"https://api.github.com/orgs/nais/issues\",\n      \"login\": \"nais\",\n      \"members_url\": \"https://api.github.com/orgs/nais/members{/member}\",\n      \"node_id\": \"MDEyOk9yZ2FuaXphdGlvbjI5NDg4Mjg5\",\n      \"public_members_url\": \"https://api.github.com/orgs/nais/public_members{/member}\",\n      \"repos_url\": \"https://api.github.com/orgs/nais/repos\",\n      \"url\": \"https://api.github.com/orgs/nais\"\n    },\n    \"pusher\": {\n      \"email\": \"38552193+johnDoe@users.noreply.github.com\",\n      \"name\": \"johnDoe\"\n    },\n    \"ref\": \"refs/heads/main\",\n    \"repository\": {\n      \"allow_forking\": true,\n      \"archive_url\": \"https://api.github.com/repos/nais/salsa/{archive_format}{/ref}\",\n      \"archived\": false,\n      \"assignees_url\": \"https://api.github.com/repos/nais/salsa/assignees{/user}\",\n      \"blobs_url\": \"https://api.github.com/repos/nais/salsa/git/blobs{/sha}\",\n      \"branches_url\": \"https://api.github.com/repos/nais/salsa/branches{/branch}\",\n      \"clone_url\": \"https://github.com/nais/salsa.git\",\n      \"collaborators_url\": \"https://api.github.com/repos/nais/salsa/collaborators{/collaborator}\",\n      \"comments_url\": \"https://api.github.com/repos/nais/salsa/comments{/number}\",\n      \"commits_url\": \"https://api.github.com/repos/nais/salsa/commits{/sha}\",\n      \"compare_url\": \"https://api.github.com/repos/nais/salsa/compare/{base}...{head}\",\n      \"contents_url\": \"https://api.github.com/repos/nais/salsa/contents/{+path}\",\n      \"contributors_url\": \"https://api.github.com/repos/nais/salsa/contributors\",\n      \"created_at\": 1639493570,\n      \"default_branch\": \"main\",\n      \"deployments_url\": \"https://api.github.com/repos/nais/salsa/deployments\",\n      \"description\": \"in line with the best from abroad\",\n      \"disabled\": false,\n      \"downloads_url\": \"https://api.github.com/repos/nais/salsa/downloads\",\n      \"events_url\": \"https://api.github.com/repos/nais/salsa/events\",\n      \"fork\": false,\n      \"forks\": 0,\n      \"forks_count\": 0,\n      \"forks_url\": \"https://api.github.com/repos/nais/salsa/forks\",\n      \"full_name\": \"nais/salsa\",\n      \"git_commits_url\": \"https://api.github.com/repos/nais/salsa/git/commits{/sha}\",\n      \"git_refs_url\": \"https://api.github.com/repos/nais/salsa/git/refs{/sha}\",\n      \"git_tags_url\": \"https://api.github.com/repos/nais/salsa/git/tags{/sha}\",\n      \"git_url\": \"git://github.com/nais/salsa.git\",\n      \"has_downloads\": true,\n      \"has_issues\": true,\n      \"has_pages\": false,\n      \"has_projects\": true,\n      \"has_wiki\": true,\n      \"homepage\": \"\",\n      \"hooks_url\": \"https://api.github.com/repos/nais/salsa/hooks\",\n      \"html_url\": \"https://github.com/nais/salsa\",\n      \"id\": 438292024,\n      \"is_template\": false,\n      \"issue_comment_url\": \"https://api.github.com/repos/nais/salsa/issues/comments{/number}\",\n      \"issue_events_url\": \"https://api.github.com/repos/nais/salsa/issues/events{/number}\",\n      \"issues_url\": \"https://api.github.com/repos/nais/salsa/issues{/number}\",\n      \"keys_url\": \"https://api.github.com/repos/nais/salsa/keys{/key_id}\",\n      \"labels_url\": \"https://api.github.com/repos/nais/salsa/labels{/name}\",\n      \"language\": \"Go\",\n      \"languages_url\": \"https://api.github.com/repos/nais/salsa/languages\",\n      \"license\": {\n        \"key\": \"mit\",\n        \"name\": \"MIT License\",\n        \"node_id\": \"MDc6TGljZW5zZTEz\",\n        \"spdx_id\": \"MIT\",\n        \"url\": \"https://api.github.com/licenses/mit\"\n      },\n      \"master_branch\": \"main\",\n      \"merges_url\": \"https://api.github.com/repos/nais/salsa/merges\",\n      \"milestones_url\": \"https://api.github.com/repos/nais/salsa/milestones{/number}\",\n      \"mirror_url\": null,\n      \"name\": \"salsa\",\n      \"node_id\": \"R_kgDOGh_OOA\",\n      \"notifications_url\": \"https://api.github.com/repos/nais/salsa/notifications{?since,all,participating}\",\n      \"open_issues\": 0,\n      \"open_issues_count\": 0,\n      \"organization\": \"nais\",\n      \"owner\": {\n        \"avatar_url\": \"https://avatars.githubusercontent.com/u/29488289?v=4\",\n        \"email\": null,\n        \"events_url\": \"https://api.github.com/users/nais/events{/privacy}\",\n        \"followers_url\": \"https://api.github.com/users/nais/followers\",\n        \"following_url\": \"https://api.github.com/users/nais/following{/other_user}\",\n        \"gists_url\": \"https://api.github.com/users/nais/gists{/gist_id}\",\n        \"gravatar_id\": \"\",\n        \"html_url\": \"https://github.com/nais\",\n        \"id\": 29488289,\n        \"login\": \"nais\",\n        \"name\": \"nais\",\n        \"node_id\": \"MDEyOk9yZ2FuaXphdGlvbjI5NDg4Mjg5\",\n        \"organizations_url\": \"https://api.github.com/users/nais/orgs\",\n        \"received_events_url\": \"https://api.github.com/users/nais/received_events\",\n        \"repos_url\": \"https://api.github.com/users/nais/repos\",\n        \"site_admin\": false,\n        \"starred_url\": \"https://api.github.com/users/nais/starred{/owner}{/repo}\",\n        \"subscriptions_url\": \"https://api.github.com/users/nais/subscriptions\",\n        \"type\": \"Organization\",\n        \"url\": \"https://api.github.com/users/nais\"\n      },\n      \"private\": false,\n      \"pulls_url\": \"https://api.github.com/repos/nais/salsa/pulls{/number}\",\n      \"pushed_at\": 1646662389,\n      \"releases_url\": \"https://api.github.com/repos/nais/salsa/releases{/id}\",\n      \"size\": 2334,\n      \"ssh_url\": \"git@github.com:nais/salsa.git\",\n      \"stargazers\": 0,\n      \"stargazers_count\": 0,\n      \"stargazers_url\": \"https://api.github.com/repos/nais/salsa/stargazers\",\n      \"statuses_url\": \"https://api.github.com/repos/nais/salsa/statuses/{sha}\",\n      \"subscribers_url\": \"https://api.github.com/repos/nais/salsa/subscribers\",\n      \"subscription_url\": \"https://api.github.com/repos/nais/salsa/subscription\",\n      \"svn_url\": \"https://github.com/nais/salsa\",\n      \"tags_url\": \"https://api.github.com/repos/nais/salsa/tags\",\n      \"teams_url\": \"https://api.github.com/repos/nais/salsa/teams\",\n      \"topics\": [\n        \"artifact\",\n        \"integrity\",\n        \"provenance\",\n        \"slsa\",\n        \"software-supply-chain\"\n      ],\n      \"trees_url\": \"https://api.github.com/repos/nais/salsa/git/trees{/sha}\",\n      \"updated_at\": \"2022-02-09T13:23:29Z\",\n      \"url\": \"https://github.com/nais/salsa\",\n      \"visibility\": \"public\",\n      \"watchers\": 0,\n      \"watchers_count\": 0\n    },\n    \"sender\": {\n      \"avatar_url\": \"https://avatars.githubusercontent.com/u/38552193?v=4\",\n      \"events_url\": \"https://api.github.com/users/johnDoe/events{/privacy}\",\n      \"followers_url\": \"https://api.github.com/users/johnDoe/followers\",\n      \"following_url\": \"https://api.github.com/users/johnDoe/following{/other_user}\",\n      \"gists_url\": \"https://api.github.com/users/johnDoe/gists{/gist_id}\",\n      \"gravatar_id\": \"\",\n      \"html_url\": \"https://github.com/johnDoe\",\n      \"id\": 38552193,\n      \"login\": \"johnDoe\",\n      \"node_id\": \"MDQ6VXNlcjM4NTUyMTkz\",\n      \"organizations_url\": \"https://api.github.com/users/johnDoe/orgs\",\n      \"received_events_url\": \"https://api.github.com/users/johnDoe/received_events\",\n      \"repos_url\": \"https://api.github.com/users/johnDoe/repos\",\n      \"site_admin\": false,\n      \"starred_url\": \"https://api.github.com/users/johnDoe/starred{/owner}{/repo}\",\n      \"subscriptions_url\": \"https://api.github.com/users/johnDoe/subscriptions\",\n      \"type\": \"User\",\n      \"url\": \"https://api.github.com/users/johnDoe\"\n    }\n  },\n  \"server_url\": \"https://github.com\",\n  \"api_url\": \"https://api.github.com\",\n  \"graphql_url\": \"https://api.github.com/graphql\",\n  \"ref_name\": \"main\",\n  \"ref_protected\": true,\n  \"ref_type\": \"branch\",\n  \"secret_source\": \"Actions\",\n  \"workspace\": \"/home/runner/work/salsa/salsa\",\n  \"action\": \"__self\",\n  \"event_path\": \"/home/runner/work/_temp/_github_workflow/event.json\",\n  \"action_repository\": \"\",\n  \"action_ref\": \"\",\n  \"path\": \"/home/runner/work/_temp/_runner_file_commands/add_path_2c66c6fd-3448-4d82-9dc0-2cea72cc64d0\",\n  \"env\": \"/home/runner/work/_temp/_runner_file_commands/set_env_2c66c6fd-3448-4d82-9dc0-2cea72cc64d0\",\n  \"step_summary\": \"/home/runner/work/_temp/_runner_file_commands/step_summary_2c66c6fd-3448-4d82-9dc0-2cea72cc64d0\"\n}",
        "INPUT_IMAGE": "ttl.sh/nais/salsa-test:1h",
        "INPUT_REPO_DIR": "/github/workspace",
        "INPUT_REPO_NAME": "nais/salsa",
        "INPUT_RUNNER_CONTEXT": "{\n  \"os\": \"Linux\",\n  \"arch\": \"X64\",\n  \"name\": \"GitHub Actions 39\",\n  \"tool_cache\": \"/opt/hostedtoolcache\",\n  \"temp\": \"/home/runner/work/_temp\",\n  \"workspace\": \"/home/runner/work/salsa\"\n}",
        "LANG": "C.UTF-8",
        "LC_COLLATE": "C",
        "PAGER": "less",
        "PS1": "\\h:\\w\\$ ",
        "PWD": "/github/workspace",
        "RUNNER_ARCH": "X64",
        "RUNNER_NAME": "GitHub Actions 39",
        "RUNNER_OS": "Linux",
        "RUNNER_TEMP": "/home/runner/work/_temp",
        "RUNNER_TOOL_CACHE": "/opt/hostedtoolcache",
        "RUNNER_WORKSPACE": "/home/runner/work/salsa",
        "SHLVL": "1"
}`

var expectedEnvs = `
{
        "CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE": "/github/workspace/gha-creds-f1d01cf9a81874ff.json",
        "CLOUDSDK_CORE_PROJECT": "plattformsikkerhet-dev-496e",
        "CLOUDSDK_PROJECT": "plattformsikkerhet-dev-496e",
        "GCLOUD_PROJECT": "plattformsikkerhet-dev-496e",
        "GCP_PROJECT": "plattformsikkerhet-dev-496e",
        "GOOGLE_APPLICATION_CREDENTIALS": "/github/workspace/gha-creds-f1d01cf9a81874ff.json",
        "GOOGLE_CLOUD_PROJECT": "plattformsikkerhet-dev-496e",
        "GOOGLE_GHA_CREDS_PATH": "/github/workspace/gha-creds-f1d01cf9a81874ff.json",
		"LC_COLLATE":"C"
}`

var envContext = `{
  		"GOVERSION": "1.17",
		"GOROOT": "/opt/hostedtoolcache/go/1.17.6/x64"
	  }`
