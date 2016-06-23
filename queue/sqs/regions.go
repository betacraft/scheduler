package sqs

// All the regions that could be used to create queues
// The map is against the Friendly names to the actual AWS Regions
var RegionNames = map[string]string{
	"USGovWest":    "us-gov-west-1",
	"USEast":       "us-east-1",
	"USWest":       "us-west-1",
	"USWest2":      "us-west-2",
	"EUWest":       "eu-west-1",
	"EUCentral":    "eu-central-1",
	"APSoutheast":  "ap-southeast-1",
	"APSoutheast2": "ap-southeast-2",
	"APNortheast":  "ap-northeast-1",
	"SAEast":       "sa-east-1",
	"CNNorth1":     "cn-north-1",
}
