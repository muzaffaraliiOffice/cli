package vars

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/services"
)

func CmdUnset(svcName, defaultSvcID, key string, iv IVars, is services.IServices) error {
	if svcName != "" {
		service, err := is.RetrieveByLabel(svcName)
		if err != nil {
			return err
		}
		if service == nil {
			return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"datica services\" command.", svcName)
		}
		defaultSvcID = service.ID
	}
	err := iv.Unset(defaultSvcID, key)
	if err != nil {
		return err
	}
	logrus.Printf("Unset. For these environment variable changes to take effect, you will need to redeploy your service with \"datica redeploy %s\"", svcName)
	return nil
}

// Unset deletes an environment variable. Any changes to environment variables
// will not take effect until the service is redeployed by pushing new code
// or via `datica redeploy`.
func (v *SVars) Unset(svcID, variable string) error {
	headers := v.Settings.HTTPManager.GetHeaders(v.Settings.SessionToken, v.Settings.Version, v.Settings.Pod, v.Settings.UsersID)
	resp, statusCode, err := v.Settings.HTTPManager.Delete(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/env/%s", v.Settings.PaasHost, v.Settings.PaasHostVersion, v.Settings.EnvironmentID, svcID, variable), headers)
	if err != nil {
		return err
	}
	return v.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
