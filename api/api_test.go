package api

import "testing"

func TestAPIManager_LoadSpecs(t *testing.T) {
	tests := []struct {
		name     string
		specPath string
		wantErr  bool
	}{
		{
			name:     "Valid spec directory",
			specPath: "../testdata/api",
			wantErr:  false,
		},
		{
			name:     "Non-existent directory",
			specPath: "../nothing_here",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewAPIManager(tt.specPath)
			err := manager.LoadSpecs()
			if (err != nil) != tt.wantErr {
				t.Errorf("APIManager.LoadSpecs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
