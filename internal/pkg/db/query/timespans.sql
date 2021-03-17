-- name: UpsertTimespansV1 :one
INSERT INTO timespans (
  schwacke_id,
  schwacke_code,
  valid_from,
  valid_until
) VALUES (
  $1, $2, $3, $4
) 
ON CONFLICT (schwacke_id) 
DO UPDATE SET 
  schwacke_code = EXCLUDED.schwacke_code,
  valid_from = EXCLUDED.valid_from,
  valid_until = EXCLUDED.valid_until
  RETURNING *;

  -- name: CountTimespansV1 :one
SELECT 
   COUNT(*) 
FROM 
   timespans;

   