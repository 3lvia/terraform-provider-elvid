terraform {
  required_version = ">= 0.13"
  required_providers {
    elvid = {
       source = "local/3lvia/elvid" # local/3lvia/elvid is used to point at an local version of the provider it will then look in $env:APPDATA\Roaming\terraform.d\plugins\local\3lvia\elvid for windows
       # source = "3lvia/elvid" # 3lvia/elvid is used to point at the published version in terraform registry. 
    }
  }
}
