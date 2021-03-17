/* table rims
code
width
height
one_piece
diameter
material
 */
CREATE TABLE "rims" (
  "code" int PRIMARY KEY,
  "width" NUMERIC(4, 2), 
  "height" char(1),
  "one_piece" boolean,
  "diameter" int,
  "material" char(1) 
);

COMMENT ON COLUMN "rims"."code" IS 'Eindeutiger Code';
COMMENT ON COLUMN "rims"."width" IS 'Felgenmaulweite n Zoll, float mit max 2 vor und 2 nachkomma stellen';
COMMENT ON COLUMN "rims"."height" IS 'Code für Felgenhöhe';
COMMENT ON COLUMN "rims"."one_piece" IS 'X = True, x = False';
COMMENT ON COLUMN "rims"."diameter" IS 'Felgendurchmesser in Zoll';
COMMENT ON COLUMN "rims"."material" IS 'S = Stahl, L = Leichtmetall';

CREATE TABLE "timespans" (
  "schwacke_id" int PRIMARY KEY,
  "schwacke_code" int,
  "valid_from" date,
  "valid_until" date
);

COMMENT ON COLUMN "timespans"."schwacke_id" IS 'interner Zeitraumschlüssel';
COMMENT ON COLUMN "timespans"."schwacke_code" IS 'Interner Typschlüssel';
COMMENT ON COLUMN "timespans"."valid_from" IS 'Datum (TT.MM.JJJJ)';
COMMENT ON COLUMN "timespans"."valid_until" IS 'Datum (TT.MM.JJJJ). Wert ist  falls gültig';

CREATE TABLE "logs" (
  "id" bigserial PRIMARY KEY,
  "inserts" int,
  "updates" int,
  "errors" int,
  "timestamp_started" timestamp,
  "timestamp_finished" timestamp
);

COMMENT ON COLUMN "logs"."inserts" IS 'Anzahl Inserts';
COMMENT ON COLUMN "logs"."updates" IS 'Anzahl Updates';
COMMENT ON COLUMN "logs"."errors" IS 'Anzahl Errors';
COMMENT ON COLUMN "logs"."timestamp_started" IS 'Startzeitpunkt';
COMMENT ON COLUMN "logs"."timestamp_finished" IS 'Endzeitpunkt';



/*
CREATE TABLE "datTables" (
  "id" bigserial PRIMARY KEY,
  "table" varchar(100),
  "version" int,
  "columns" int,
);

CREATE TABLE "datColums" (
  "id" bigserial PRIMARY KEY,
  "columnName" varchar(100),
  "dataType" int,
  "required" boolean,
  "length" int,
  "minValueType" 
);
*/