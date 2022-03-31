package location

import (
	"testing"

	"github.com/RafayLabs/rcloud-cli/pkg/config"
	"github.com/RafayLabs/rcloud-cli/pkg/context"
	// "github.com/RafayLabs/rcloud-cli/pkg/models"
)

func TestListLocations(t *testing.T) {
	// set up the config
	t.Log("Setting up config")
	cliCtx := context.GetContext()
	cliCtx.ConfigDir = "../../testdata/env/"
	err := config.InitConfig(cliCtx)
	if err != nil {
		t.Errorf("Error setting up config: %v", err)
	}

	type args struct {
		partner string
		limit   int
		offset  int
	}
	tests := []struct {
		name    string
		args    args
		// want    []*models.Metro
		want1   int
		wantErr bool
	}{
		{
			name: "Test list locations with limit",
			args: args{
				limit:  1000,
				offset: 0,
			},
			wantErr: false,
		},
		{
			name: "Test list locations with limit and offset",
			args: args{
				limit:  1,
				offset: 1,
			},
			wantErr: false,
		},
		{
			name: "Test list locations with negative limit",
			args: args{
				limit:  -1000,
				offset: 0,
			},
			wantErr: true,
		},
		{
			name: "Test list locations with negative offset",
			args: args{
				limit:  1000,
				offset: -1,
			},
			wantErr: true,
		},
		{
			name: "Test list locations with negative limit and offset",
			args: args{
				limit:  -1000,
				offset: -1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := ListLocation(tt.args.partner, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListLocatin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
