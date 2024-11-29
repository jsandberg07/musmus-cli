-- settings is saving the settings. this should ONLY ever have one row

-- name: SetUpSettings :exec
INSERT INTO settings(settings_complete, only_activate_self)
VALUES (false, false);

-- name: GetSettings :one
SELECT * FROM settings;

-- name: ActivateSelfOnly :one
SELECT only_activate_self FROM settings;

-- name: UpdateActivateSelf :exec
UPDATE settings SET only_activate_self = $1;

-- name: FirstTimeSetupComplete :exec
UPDATE settings SET settings_complete = true;