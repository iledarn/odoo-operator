package stub

import (
	"context"

	api "github.com/xoes/odoo-operator/pkg/apis/odoo/v1alpha1"
	"github.com/xoes/odoo-operator/pkg/cluster"
	"github.com/xoes/odoo-operator/pkg/instance"
	"github.com/xoes/odoo-operator/pkg/pg"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *api.OdooCluster:
		return cluster.Reconcile(o)
	case *api.PgNamespace:
		return pg.ReconcilePgNamespace(o)
	case *api.ClusterMigration:
		return cluster.ReconcileMigration(o)
	case *api.OdooInstance:
		return instance.Reconcile(o)
	case *api.InstanceMigration:
		return instance.CreateMigration(o)
	case *api.InstanceBackup:
		return instance.CreateBackup(o)
	}
	return nil
}
