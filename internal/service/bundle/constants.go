// Package bundle contains constants and assembly logic for the v2 snapshot publisher.
// Constants are ported verbatim from indoormap-backend/indoormap_api/utils/constants.py.
package bundle

// ExportableMapping maps exhibitor_english_status int to [en, ja] label pair.
// Django: FOODEX_EXPORTABLE_MAPPING / HCJ_EXPORTABLE_MAPPING
type ExportableEntry [2]string // [en, ja]

// ExportableMapping is keyed by exhibitor_english_status value.
type ExportableMapping map[int]ExportableEntry

// ImporterMapping maps importer option int to [en, ja] label pair.
type ImporterEntry [2]string // [en, ja]
type ImporterMapping map[int]ImporterEntry

// RequisiteDocument is [en, ja] label pair.
type RequisiteDocument [2]string // [en, ja]

// FOODEX venue constants (venue.ExternalID contains "foodex")
var FoodexExportableMapping = ExportableMapping{
	1: {"Exportable", "輸出可能"},
	2: {"Not exportable", "輸出不可"},
}

var FoodexImporterMapping = ImporterMapping{
	1: {"Already decided", "インポーターが決定している"},
	2: {"Not yet decided", "インポーターが決定していない"},
}

// HCJ venue constants (venue.ExternalID contains "hcj")
var HCJExportableMapping = ExportableMapping{
	1: {"Available", "海外取引可能"},
	2: {"Not available", "海外取引不可"},
}

var HCJImporterMapping = ImporterMapping{
	1: {"Already decided", "インポーターが決定している"},
	2: {"Not yet decided", "インポーターが決定していない"},
}

// ProductRequisiteDocuments ported from indoormap_api/utils/constants.py:133.
var ProductRequisiteDocuments = []RequisiteDocument{
	{"ISO9001", "ISO9001"},
	{"ISO22000", "ISO22000"},
	{"EU-HACCP", "EU-HACCP"},
	{"US-HACCP", "米国HACCP"},
	{"Other HACCP", "その他HACCP"},
	{"HALAL", "HALAL"},
	{"FSMA", "FSMA"},
	{"FDA", "FDA"},
	{"GMP", "GMP"},
	{"GAP", "GAP"},
	{"CE Mark", "CEマーク"},
	{"CCC(China Compulsory Certification)", "CCC（中国強制認証制度）"},
	{"UKCA(UK Conformity Assessed) marking", "UKCA（英国適合性評価）"},
	{"Organic JAS Mark", "有機JASマーク"},
	{"USDA Organic Certification", "USDAオーガニック認証"},
	{"EU Organic Certification", "EUオーガニック認証"},
	{"Non-GMO Certification", "非遺伝子組み換え認証"},
	{"GFCO(Gluten-Free Certification Organization)", "GFCO"},
	{"GFCP(Gluten Free Certification Program)", "GFCP"},
	{"MSC(Marine Stewardship Council)", "MSC"},
	{"ASC(Aquaculture Stewardship Council)", "ASC"},
	{"Marine Eco-Label Japan Council", "マリンエコラベルジャパン"},
	{"Aquaculture Eco-Label(AEL)", "養殖エコラベル"},
	{"Others", "その他"},
}

// SelectMapping returns the FOODEX or HCJ mapping based on the venue external ID.
func SelectExportableMapping(venueExternalID string) ExportableMapping {
	if contains(venueExternalID, "hcj") {
		return HCJExportableMapping
	}
	return FoodexExportableMapping
}

func SelectImporterMapping(venueExternalID string) ImporterMapping {
	if contains(venueExternalID, "hcj") {
		return HCJImporterMapping
	}
	return FoodexImporterMapping
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
