package search

// VenueIndexMapping returns the complete OpenSearch index mapping for a venue index.
// It matches the Python get_unified_mapping() function exactly.
//
// Required OpenSearch plugins on the cluster:
//   - analysis-kuromoji  (Japanese morphological analysis)
//   - analysis-icu       (Unicode normalization and collation)
//   - analysis-phonetic  (Double Metaphone for typo tolerance)
func VenueIndexMapping() map[string]any {
	return map[string]any{
		"settings": venueIndexSettings(),
		"mappings": map[string]any{
			"properties": venueIndexProperties(),
		},
	}
}

func venueIndexSettings() map[string]any {
	return map[string]any{
		"number_of_shards":   3,
		"number_of_replicas": 2,
		"max_result_window":  20000,
		"refresh_interval":   "10s",
		"analysis": map[string]any{
			"char_filter": map[string]any{
				"normalize_chars": map[string]any{
					"type": "mapping",
					"mappings": []string{
						"０=>0", "１=>1", "２=>2", "３=>3", "４=>4",
						"５=>5", "６=>6", "７=>7", "８=>8", "９=>9",
						"（=>(", "）=>)", "　=> ", "－=>-", "ー=>-",
						"～=>~", "・=> ",
						"Ａ=>A", "Ｂ=>B", "Ｃ=>C", "Ｄ=>D", "Ｅ=>E",
						"Ｆ=>F", "Ｇ=>G", "Ｈ=>H", "Ｉ=>I", "Ｊ=>J",
						"Ｋ=>K", "Ｌ=>L", "Ｍ=>M", "Ｎ=>N", "Ｏ=>O",
						"Ｐ=>P", "Ｑ=>Q", "Ｒ=>R", "Ｓ=>S", "Ｔ=>T",
						"Ｕ=>U", "Ｖ=>V", "Ｗ=>W", "Ｘ=>X", "Ｙ=>Y", "Ｚ=>Z",
						"ａ=>a", "ｂ=>b", "ｃ=>c", "ｄ=>d", "ｅ=>e",
						"ｆ=>f", "ｇ=>g", "ｈ=>h", "ｉ=>i", "ｊ=>j",
						"ｋ=>k", "ｌ=>l", "ｍ=>m", "ｎ=>n", "ｏ=>o",
						"ｐ=>p", "ｑ=>q", "ｒ=>r", "ｓ=>s", "ｔ=>t",
						"ｕ=>u", "ｖ=>v", "ｗ=>w", "ｘ=>x", "ｙ=>y", "ｚ=>z",
						"！=>!", "？=>?", "＠=>@",
					},
				},
				"katakana_normalize": map[string]any{
					"type": "mapping",
					"mappings": []string{
						// Small kana to full kana
						"ァ=>ア", "ィ=>イ", "ゥ=>ウ", "ェ=>エ", "ォ=>オ",
						"ャ=>ヤ", "ュ=>ユ", "ョ=>ヨ", "ッ=>ツ", "ヮ=>ワ",
						"ヵ=>カ", "ヶ=>ケ",
						// Halfwidth katakana to fullwidth
						"ｱ=>ア", "ｲ=>イ", "ｳ=>ウ", "ｴ=>エ", "ｵ=>オ",
						"ｶ=>カ", "ｷ=>キ", "ｸ=>ク", "ｹ=>ケ", "ｺ=>コ",
						"ｻ=>サ", "ｼ=>シ", "ｽ=>ス", "ｾ=>セ", "ｿ=>ソ",
						"ﾀ=>タ", "ﾁ=>チ", "ﾂ=>ツ", "ﾃ=>テ", "ﾄ=>ト",
						"ﾅ=>ナ", "ﾆ=>ニ", "ﾇ=>ヌ", "ﾈ=>ネ", "ﾉ=>ノ",
						"ﾊ=>ハ", "ﾋ=>ヒ", "ﾌ=>フ", "ﾍ=>ヘ", "ﾎ=>ホ",
						"ﾏ=>マ", "ﾐ=>ミ", "ﾑ=>ム", "ﾒ=>メ", "ﾓ=>モ",
						"ﾔ=>ヤ", "ﾕ=>ユ", "ﾖ=>ヨ",
						"ﾗ=>ラ", "ﾘ=>リ", "ﾙ=>ル", "ﾚ=>レ", "ﾛ=>ロ",
						"ﾜ=>ワ", "ｦ=>ヲ", "ﾝ=>ン",
						// Small halfwidth katakana
						"ｧ=>ア", "ｨ=>イ", "ｩ=>ウ", "ｪ=>エ", "ｫ=>オ",
						"ｬ=>ヤ", "ｭ=>ユ", "ｮ=>ヨ", "ｯ=>ツ",
					},
				},
			},
			"tokenizer": map[string]any{
				"kuromoji_tokenizer": map[string]any{
					"type":                 "kuromoji_tokenizer",
					"mode":                 "search",
					"discard_punctuation":  false,
				},
			},
			"filter": map[string]any{
				"japanese_stop": map[string]any{
					"type": "stop",
					"stopwords": []string{
						"の", "に", "は", "を", "が", "で", "と", "から",
						"まで", "も", "や", "など", "へ", "より", "まで",
						"として", "こと", "よう", "する", "いる", "れる",
						"られる", "です", "ます", "など", "おり",
					},
				},
				"english_stop": map[string]any{
					"type":      "stop",
					"stopwords": "_english_",
				},
				"katakana_stem": map[string]any{
					"type":           "kuromoji_stemmer",
					"minimum_length": 2,
				},
				"edge_ngram": map[string]any{
					"type":     "edge_ngram",
					"min_gram": 2,
					"max_gram": 15,
				},
				"ngram_filter": map[string]any{
					"type":     "ngram",
					"min_gram": 3,
					"max_gram": 3,
				},
				"english_stemmer": map[string]any{
					"type":     "stemmer",
					"language": "english",
				},
				"kuromoji_readingform": map[string]any{
					"type":        "kuromoji_readingform",
					"use_romaji":  false,
				},
				"icu_normalizer": map[string]any{
					"type": "icu_normalizer",
					"name": "nfkc",
				},
				"synonym_filter": map[string]any{
					"type":      "synonym",
					"tokenizer": "kuromoji_tokenizer",
					"synonyms": []string{
						"コーヒー,Coffee",
						"ワイン,Wine",
						"ビール,Beer",
						"チーズ,Cheese",
						"パン,Bread",
						"飲料,Beverages",
						"食品,Food",
					},
				},
				"cjk_bigram_all": map[string]any{
					"type":            "cjk_bigram",
					"output_unigrams": true,
				},
				"search_edge_ngram": map[string]any{
					"type":     "edge_ngram",
					"min_gram": 3,
					"max_gram": 10,
				},
				"fuzzy_ngram": map[string]any{
					"type":     "edge_ngram",
					"min_gram": 2,
					"max_gram": 4,
				},
				"phonetic_filter": map[string]any{
					"type":    "phonetic",
					"encoder": "double_metaphone",
					"replace": false,
				},
			},
			"analyzer": map[string]any{
				"japanese_analyzer": map[string]any{
					"char_filter": []string{"normalize_chars", "katakana_normalize"},
					"tokenizer":   "kuromoji_tokenizer",
					"filter": []string{
						"kuromoji_baseform", "kuromoji_part_of_speech", "kuromoji_readingform",
						"katakana_stem", "cjk_width", "icu_normalizer", "lowercase",
						"japanese_stop", "synonym_filter", "cjk_bigram_all", "ngram_filter",
					},
				},
				"japanese_search": map[string]any{
					"char_filter": []string{"normalize_chars", "katakana_normalize"},
					"tokenizer":   "kuromoji_tokenizer",
					"filter": []string{
						"kuromoji_baseform", "kuromoji_part_of_speech", "kuromoji_readingform",
						"katakana_stem", "cjk_width", "icu_normalizer", "lowercase",
						"japanese_stop", "synonym_filter",
					},
				},
				"japanese_autocomplete": map[string]any{
					"char_filter": []string{"normalize_chars", "katakana_normalize"},
					"tokenizer":   "kuromoji_tokenizer",
					"filter": []string{
						"kuromoji_baseform", "kuromoji_part_of_speech", "kuromoji_readingform",
						"katakana_stem", "cjk_width", "icu_normalizer", "lowercase",
						"japanese_stop", "edge_ngram",
					},
				},
				"japanese_fuzzy": map[string]any{
					"char_filter": []string{"normalize_chars", "katakana_normalize"},
					"tokenizer":   "kuromoji_tokenizer",
					"filter": []string{
						"kuromoji_baseform", "kuromoji_part_of_speech", "kuromoji_readingform",
						"katakana_stem", "cjk_width", "icu_normalizer", "lowercase",
						"japanese_stop", "search_edge_ngram",
					},
				},
				"english_analyzer": map[string]any{
					"tokenizer": "standard",
					"filter": []string{
						"lowercase", "asciifolding", "english_stop",
						"synonym_filter", "english_stemmer", "ngram_filter",
					},
				},
				"english_search": map[string]any{
					"tokenizer": "standard",
					"filter":    []string{"lowercase", "asciifolding", "synonym_filter"},
				},
				"english_autocomplete": map[string]any{
					"tokenizer": "standard",
					"filter":    []string{"lowercase", "asciifolding", "edge_ngram"},
				},
				"fuzzy_search": map[string]any{
					"tokenizer": "standard",
					"filter":    []string{"lowercase", "asciifolding", "search_edge_ngram"},
				},
				"typo_tolerant": map[string]any{
					"tokenizer": "standard",
					"filter":    []string{"lowercase", "asciifolding", "phonetic_filter", "fuzzy_ngram"},
				},
			},
		},
	}
}

func venueIndexProperties() map[string]any {
	return map[string]any{
		"type": map[string]any{
			"type": "keyword",
		},
		"venue": map[string]any{
			"properties": map[string]any{
				"external_id": map[string]any{"type": "keyword"},
				"section":     map[string]any{"type": "integer"},
			},
		},
		"search_text": map[string]any{
			"properties": map[string]any{
				"ja": map[string]any{
					"type":            "text",
					"analyzer":        "japanese_analyzer",
					"search_analyzer": "japanese_search",
					"fields": map[string]any{
						"autocomplete": map[string]any{
							"type":            "text",
							"analyzer":        "japanese_autocomplete",
							"search_analyzer": "japanese_search",
						},
						"keyword": map[string]any{"type": "keyword"},
					},
				},
				"en": map[string]any{
					"type":            "text",
					"analyzer":        "english_analyzer",
					"search_analyzer": "english_search",
					"fields": map[string]any{
						"autocomplete": map[string]any{
							"type":            "text",
							"analyzer":        "english_autocomplete",
							"search_analyzer": "english_search",
						},
						"keyword": map[string]any{"type": "keyword"},
					},
				},
				"all": map[string]any{
					"type":            "text",
					"analyzer":        "japanese_analyzer",
					"search_analyzer": "japanese_search",
				},
			},
		},
		"product": map[string]any{
			"properties": map[string]any{
				"id": map[string]any{"type": "keyword"},
				"name": map[string]any{
					"type":            "text",
					"analyzer":        "japanese_analyzer",
					"search_analyzer": "japanese_search",
					"copy_to":         "search_text.all",
					"fields": map[string]any{
						"keyword":      map[string]any{"type": "keyword"},
						"autocomplete": map[string]any{"type": "text", "analyzer": "japanese_autocomplete", "search_analyzer": "japanese_search"},
						"fuzzy":        map[string]any{"type": "text", "analyzer": "japanese_fuzzy", "search_analyzer": "japanese_search"},
					},
				},
				"name_en": map[string]any{
					"type":            "text",
					"analyzer":        "english_analyzer",
					"search_analyzer": "english_search",
					"copy_to":         "search_text.all",
					"fields": map[string]any{
						"keyword":      map[string]any{"type": "keyword"},
						"autocomplete": map[string]any{"type": "text", "analyzer": "english_autocomplete", "search_analyzer": "english_search"},
						"fuzzy":        map[string]any{"type": "text", "analyzer": "fuzzy_search", "search_analyzer": "english_search"},
					},
				},
				"name_kana": map[string]any{
					"type":     "text",
					"analyzer": "japanese_analyzer",
					"fields":   map[string]any{"keyword": map[string]any{"type": "keyword"}},
				},
				"category": map[string]any{
					"type":     "text",
					"analyzer": "japanese_analyzer",
					"fields":   map[string]any{"keyword": map[string]any{"type": "keyword"}},
				},
				"category_en": map[string]any{
					"type":     "text",
					"analyzer": "english_analyzer",
					"fields":   map[string]any{"keyword": map[string]any{"type": "keyword"}},
				},
				"keywords": map[string]any{
					"type":            "text",
					"analyzer":        "japanese_analyzer",
					"search_analyzer": "japanese_search",
					"copy_to":         "search_text.all",
				},
				"description": map[string]any{
					"type":            "text",
					"analyzer":        "japanese_analyzer",
					"search_analyzer": "japanese_search",
					"copy_to":         "search_text.all",
				},
				"price": map[string]any{
					"type":   "text",
					"fields": map[string]any{"keyword": map[string]any{"type": "keyword"}},
				},
				"ingredients": map[string]any{
					"type":            "text",
					"analyzer":        "japanese_analyzer",
					"search_analyzer": "japanese_search",
				},
				"target": map[string]any{
					"type":            "text",
					"analyzer":        "japanese_analyzer",
					"search_analyzer": "japanese_search",
				},
				"use": map[string]any{
					"type":            "text",
					"analyzer":        "japanese_analyzer",
					"search_analyzer": "japanese_search",
				},
			},
		},
		"exhibitor": map[string]any{
			"properties": map[string]any{
				"id": map[string]any{"type": "keyword"},
				"name": map[string]any{
					"type":            "text",
					"analyzer":        "japanese_analyzer",
					"search_analyzer": "japanese_search",
					"copy_to":         "search_text.all",
					"fields": map[string]any{
						"keyword":      map[string]any{"type": "keyword"},
						"autocomplete": map[string]any{"type": "text", "analyzer": "japanese_autocomplete", "search_analyzer": "japanese_search"},
						"fuzzy":        map[string]any{"type": "text", "analyzer": "japanese_fuzzy", "search_analyzer": "japanese_search"},
					},
				},
				"name_en": map[string]any{
					"type":            "text",
					"analyzer":        "english_analyzer",
					"search_analyzer": "english_search",
					"copy_to":         "search_text.all",
					"fields": map[string]any{
						"keyword":      map[string]any{"type": "keyword"},
						"autocomplete": map[string]any{"type": "text", "analyzer": "english_autocomplete", "search_analyzer": "english_search"},
						"fuzzy":        map[string]any{"type": "text", "analyzer": "fuzzy_search", "search_analyzer": "english_search"},
					},
				},
				"name_kana": map[string]any{
					"type":     "text",
					"analyzer": "japanese_analyzer",
					"fields":   map[string]any{"keyword": map[string]any{"type": "keyword"}},
				},
				"booth_number": map[string]any{"type": "keyword"},
				"exhibition_zone": map[string]any{
					"type":            "text",
					"analyzer":        "japanese_analyzer",
					"copy_to":         "search_text.all",
					"fields": map[string]any{
						"keyword": map[string]any{"type": "keyword"},
						"fuzzy":   map[string]any{"type": "text", "analyzer": "japanese_fuzzy", "search_analyzer": "japanese_search"},
					},
				},
				"exhibition_zone_en": map[string]any{
					"type":     "text",
					"analyzer": "english_analyzer",
					"copy_to":  "search_text.all",
					"fields": map[string]any{
						"keyword": map[string]any{"type": "keyword"},
						"fuzzy":   map[string]any{"type": "text", "analyzer": "fuzzy_search", "search_analyzer": "english_search"},
					},
				},
				"country": map[string]any{
					"type":     "text",
					"analyzer": "japanese_analyzer",
					"fields":   map[string]any{"keyword": map[string]any{"type": "keyword"}},
				},
				"country_en": map[string]any{
					"type":     "text",
					"analyzer": "english_analyzer",
					"fields":   map[string]any{"keyword": map[string]any{"type": "keyword"}},
				},
				"description": map[string]any{
					"type":            "text",
					"analyzer":        "japanese_analyzer",
					"search_analyzer": "japanese_search",
					"copy_to":         "search_text.all",
				},
				"description_en": map[string]any{
					"type":            "text",
					"analyzer":        "english_analyzer",
					"search_analyzer": "english_search",
					"copy_to":         "search_text.all",
				},
				"products_summary": map[string]any{
					"type":            "text",
					"analyzer":        "japanese_analyzer",
					"search_analyzer": "japanese_search",
					"copy_to":         "search_text.all",
				},
				"products_summary_en": map[string]any{
					"type":            "text",
					"analyzer":        "english_analyzer",
					"search_analyzer": "english_search",
					"copy_to":         "search_text.all",
				},
				"main_logo":       map[string]any{"type": "keyword", "index": false, "norms": false},
				"top_logo":        map[string]any{"type": "keyword", "index": false, "norms": false},
				"exhibitor_image": map[string]any{"type": "keyword", "index": false, "norms": false},
				"main_image":      map[string]any{"type": "keyword", "index": false, "norms": false},
				"custom_images":   map[string]any{"type": "keyword", "index": false, "norms": false},
				"file_list": map[string]any{
					"type": "nested",
					"properties": map[string]any{
						"file":  map[string]any{"type": "keyword", "index": false, "norms": false},
						"title": map[string]any{"type": "text", "analyzer": "japanese_analyzer", "index": false},
					},
				},
				"videos": map[string]any{"type": "keyword", "index": false, "norms": false},
				"attachments": map[string]any{
					"type": "nested",
					"properties": map[string]any{
						"title":      map[string]any{"type": "text", "analyzer": "japanese_analyzer", "index": false},
						"file_type":  map[string]any{"type": "keyword", "index": false, "norms": false},
						"source_url": map[string]any{"type": "keyword", "index": false, "norms": false},
						"file_url":   map[string]any{"type": "keyword", "index": false, "norms": false},
					},
				},
				"files": map[string]any{
					"type": "nested",
					"properties": map[string]any{
						"title":      map[string]any{"type": "text", "analyzer": "japanese_analyzer", "index": false},
						"file_type":  map[string]any{"type": "keyword", "index": false, "norms": false},
						"source_url": map[string]any{"type": "keyword", "index": false, "norms": false},
						"file_url":   map[string]any{"type": "keyword", "index": false, "norms": false},
					},
				},
			},
		},
		"categories": map[string]any{
			"type": "nested",
			"properties": map[string]any{
				"id":      map[string]any{"type": "keyword"},
				"name":    map[string]any{"type": "text", "analyzer": "japanese_analyzer", "fields": map[string]any{"keyword": map[string]any{"type": "keyword"}}},
				"name_en": map[string]any{"type": "text", "analyzer": "english_analyzer", "fields": map[string]any{"keyword": map[string]any{"type": "keyword"}}},
				"is_main": map[string]any{"type": "boolean"},
			},
		},
		"indexed_at": map[string]any{
			"type":   "date",
			"format": "strict_date_optional_time",
		},
	}
}
