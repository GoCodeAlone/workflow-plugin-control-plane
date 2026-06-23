package actions_test

import (
	"strings"
	"testing"

	"github.com/GoCodeAlone/workflow-plugin-control-plane/actions"
)

func TestProviderActionValidatesNeutralHandoffShape(t *testing.T) {
	action := validProviderAction()
	if err := action.Validate(); err != nil {
		t.Fatalf("provider action invalid: %v", err)
	}
}

func TestProviderActionRejectsInvalidData(t *testing.T) {
	cases := map[string]func(*actions.ProviderAction){
		"missing plugin id": func(a *actions.ProviderAction) {
			a.PluginID = ""
		},
		"control whitespace in contract": func(a *actions.ProviderAction) {
			a.Contract = "workflow-plugin-payments:step.reserve\nbad"
		},
		"missing config value": func(a *actions.ProviderAction) {
			a.Config[0].Value = ""
			a.Config[0].Values = nil
		},
		"nested too deeply": func(a *actions.ProviderAction) {
			a.Fallbacks[0].Fallbacks = []actions.ProviderAction{{
				PluginID: "workflow-plugin-payments",
				StepType: "step.reserve_fallback_nested",
				Contract: "workflow-plugin-payments:step.reserve_fallback_nested",
			}}
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			action := validProviderAction()
			mutate(&action)
			if err := action.Validate(); err == nil {
				t.Fatal("expected provider action to fail")
			}
		})
	}
}

func TestProviderActionReportsMultipleErrors(t *testing.T) {
	action := validProviderAction()
	action.PluginID = ""
	action.RequiredConfig[0] = "bad\nfield"
	action.Config[0].Name = ""
	err := action.Validate()
	if err == nil {
		t.Fatal("expected provider action errors")
	}
	for _, want := range []string{"plugin_id", "required_config[0]", "config[0]"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error %q missing %q", err, want)
		}
	}
}

func validProviderAction() actions.ProviderAction {
	return actions.ProviderAction{
		PluginID:        "workflow-plugin-payments",
		StepType:        "step.settlement_wallet_reserve",
		Contract:        "workflow-plugin-payments:step.settlement_wallet_reserve",
		ContractVersion: "v1",
		RequiredConfig:  []string{"module", "settlement_wallet_ref"},
		Config: []actions.ProviderConfig{{
			Name:  "settlement_wallet_ref",
			Value: "wallet://settlement/route-bmw-1/window-1",
		}, {
			Name:   "networks",
			Values: []string{"base", "ethereum"},
		}},
		OutputHandling: []string{"submit verified evidence without custody credentials"},
		SubmitEndpoint: "/v1/settlement/wallets/wallet-window-1/reserve",
		Prerequisites: []actions.ProviderAction{{
			PluginID:       "workflow-plugin-payments",
			StepType:       "step.webhook_endpoint_ensure",
			Contract:       "workflow-plugin-payments:step.webhook_endpoint_ensure",
			RequiredConfig: []string{"module", "url", "signing_secret_sink"},
			Config: []actions.ProviderConfig{{
				Name:   "events",
				Values: []string{"payment_intent.succeeded"},
			}},
			OutputHandling: []string{"store signing secrets outside the compute host"},
		}},
		Fallbacks: []actions.ProviderAction{{
			PluginID:       "workflow-plugin-payments",
			StepType:       "step.settlement_wallet_freeze",
			Contract:       "workflow-plugin-payments:step.settlement_wallet_freeze",
			RequiredConfig: []string{"module", "settlement_wallet_ref", "reason"},
			OutputHandling: []string{"submit freeze reason only"},
			SubmitEndpoint: "/v1/settlement/wallets/wallet-window-1/freeze",
		}},
	}
}
