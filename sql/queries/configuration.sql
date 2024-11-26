-- config is saving the settings. this should ONLY ever have one row

-- name: SetUpConfig :exec
INSERT INTO config(config_complete, only_activate_self)
VALUES (false, false);

-- name: GetConfig :one
SELECT * FROM config;

-- name: ActivateSelfOnly :one
SELECT only_activate_self FROM config;

-- name: UpdateActivateSelf :exec
UPDATE config SET only_activate_self = $1;

-- name: FirstTimeSetupComplete :exec
UPDATE config SET config_complete = true;