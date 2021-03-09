
-- name: UpsertRimsV1 :one
INSERT INTO rims (
  code,
  width,
  height,
  one_piece,
  diameter,
  material
) VALUES (
  $1, $2, $3, $4, $5, $6
) 
ON CONFLICT (code) 
DO UPDATE SET 
  width = EXCLUDED.width,
  height = EXCLUDED.height,
  one_piece = EXCLUDED.one_piece,
  diameter = EXCLUDED.diameter,
  material = EXCLUDED.material
RETURNING *;

-- name: CountRimsV1 :one
SELECT 
   COUNT(*) 
FROM 
   rims;

   