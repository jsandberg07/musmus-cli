-- settings is saving the settings. this should ONLY ever have one row

-- name: SetUpSettings :exec
INSERT INTO settings(id, settings_complete, only_activate_self, test_data_loaded)
VALUES (1, false, false, false);

-- name: GetSettings :one
SELECT * FROM settings
WHERE id = 1;

-- name: ActivateSelfOnly :one
SELECT only_activate_self FROM settings;

-- name: UpdateActivateSelf :exec
UPDATE settings SET only_activate_self = $1;

-- name: FirstTimeSetupComplete :exec
UPDATE settings SET settings_complete = true;

-- name: TestDataLoaded :exec
UPDATE settings SET test_data_loaded = true;

-- name: UpdateSettings :exec
UPDATE settings
SET only_activate_self = $1;