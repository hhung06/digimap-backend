-- QR codes and amenities are dropped features: not ported from Django, and the
-- editor's remaining calls go to the legacy backend. Remove the unused tables.
DROP TABLE IF EXISTS qrcodes;
DROP TABLE IF EXISTS venue_amenities;
DROP TABLE IF EXISTS amenities;
