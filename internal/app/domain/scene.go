package domain

import "context"

// Scene 场景
type Scene string

const (
	List   Scene = "list"
	Detail Scene = "detail"
	Create Scene = "create"
	Update Scene = "update"
	Delete Scene = "delete"
)

var sceneContextKey = Scene("scene")

func SetSceneWithContext(ctx context.Context, value Scene) context.Context {
	return context.WithValue(ctx, sceneContextKey, value)
}

func GetSceneFromContext(ctx context.Context) Scene {
	value := ctx.Value(sceneContextKey)
	if value != nil {
		return value.(Scene)
	}
	return ""
}
