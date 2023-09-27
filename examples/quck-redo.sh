cd /home/olytvyn/go/terraform-provider-firehydrant/
make install
cd /home/olytvyn/go/terraform-provider-firehydrant/examples
rm /home/olytvyn/go/terraform-provider-firehydrant/examples/.terraform.lock.hcl
terraform init -plugin-dir="/home/olytvyn/.terraform.d/plugins"
terraform plan
