@Library(value="service-builder-shared-library@master", changelog=false) _
import com.karhoo.Constants

CICD {
  helmCharts = ["core-braintree"]
  mysqlRequired = true
  containerImages = [:]
  containerImages["builder"] = [name:"karhoo-golang-mysql",tag:"0.6.5"]
  apiTestsToRun = ["v1-bookings-follow-code", "v2-payments"]
  scratchConf = "scratch-backendCICD.yaml"
  makeTargets = ["deps", "test", "dbtest", "lint", "coveralls", "build"]
  envVars = [
    'GO111MODULE': 'on',
    'GOPROXY': 'http://athens',
  ]
  vaultEnvVars = [
    [key: "BRAINTREE_PRIVATE_KEY", path: "kubernetes/scratch/_all/common", field: "braintree-private-key"]
  ]
}
