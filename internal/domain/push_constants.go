// Code generated from indoormap-backend/indoormap_api/utils/constants.py; DO NOT EDIT.
package domain

// PUSH_TYPES and TOPIC_* are ported verbatim from the Django constants file.
// They stay as generic map data so callers can consume the exact original shape.

var PUSH_TYPES = []map[string]any{
	map[string]any{
		"id":       1,
		"key":      "topic_all_users",
		"label":    "アプリ全ユーザー（TOPIC）",
		"label_en": "All app users (TOPIC)",
		"type":     nil,
	},
	map[string]any{
		"id":       2,
		"key":      "topic_exhibitor_staff_all",
		"label":    "出展者スタッフ全員（TOPIC）",
		"label_en": "All exhibitor staff (TOPIC)",
		"type":     nil,
	},
	map[string]any{
		"id":       3,
		"key":      "topic_visitor_all",
		"label":    "来場者全員（TOPIC）",
		"label_en": "All visitors (TOPIC)",
		"type":     nil,
	},
	map[string]any{
		"id":       4,
		"key":      "seminar_registered",
		"label":    "セミナー聴講登録者",
		"label_en": "Registered seminar attendees",
		"type":     "number",
		"options":  []map[string]any{},
	},
	map[string]any{
		"id":       7,
		"key":      "checked_in_visitors",
		"label":    "日チェックイン来場者",
		"label_en": "Visitors checked in (by date)",
		"type":     "date",
	},
	map[string]any{
		"id":       10,
		"key":      "topic_exhibitor_domestic",
		"label":    "国内出展者（TOPIC）",
		"label_en": "Domestic exhibitors (TOPIC)",
		"type":     nil,
	},
	map[string]any{
		"id":       11,
		"key":      "topic_exhibitor_oversea",
		"label":    "海外出展者（TOPIC）",
		"label_en": "Overseas exhibitors (TOPIC)",
		"type":     nil,
	},
	map[string]any{
		"id":       12,
		"key":      "topic_visitor_domestic",
		"label":    "国内来場者（TOPIC）",
		"label_en": "Domestic visitors (TOPIC)",
		"type":     nil,
	},
	map[string]any{
		"id":       13,
		"key":      "topic_visitor_oversea",
		"label":    "海外来場者（TOPIC）",
		"label_en": "Overseas visitors (TOPIC)",
		"type":     nil,
	},
	map[string]any{
		"id":       14,
		"key":      "topic_visitor_industry",
		"label":    "来場者業種別（TOPIC）",
		"label_en": "Visitors by industry (TOPIC)",
		"type":     "number",
		"options":  []map[string]any{},
	},
	map[string]any{
		"id":       15,
		"key":      "topic_visitor_job",
		"label":    "来場者職種別（TOPIC）",
		"label_en": "Visitors by job type (TOPIC)",
		"type":     "number",
		"options":  []map[string]any{},
	},
	map[string]any{
		"id":       16,
		"key":      "topic_visitor_position",
		"label":    "来場者役職別（TOPIC）",
		"label_en": "Visitors by job position (TOPIC)",
		"type":     "number",
		"options":  []map[string]any{},
	},
	map[string]any{
		"id":       17,
		"key":      "topic_visitor_interest",
		"label":    "来場者関心カテゴリー（TOPIC）",
		"label_en": "Visitors by interest category (TOPIC)",
		"type":     "number",
		"options":  []map[string]any{},
	},
	map[string]any{
		"id":       18,
		"key":      "topic_exhibitor_zone",
		"label":    "出展ゾーン別スタッフ（TOPIC）",
		"label_en": "Exhibitor staff by zone (TOPIC)",
		"type":     "number",
		"options":  []map[string]any{},
	},
	map[string]any{
		"id":       19,
		"key":      "topic_visitor_survey_segment_purchase",
		"label":    "【アンケート】商品の仕入・購買への関与度　（TOPIC）",
		"label_en": "[Survey] Involvement in Product Procurement / Purchasing",
		"type":     "number",
		"options":  []map[string]any{},
	},
	map[string]any{
		"id":       20,
		"key":      "topic_visitor_survey_segment_purpose",
		"label":    "【アンケート】来場目的　（TOPIC）",
		"label_en": "[Survey] Purpose of Visit　（TOPIC）",
		"type":     "number",
		"options":  []map[string]any{},
	},
	map[string]any{
		"id":       21,
		"key":      "topic_visitor_survey_segment_budget",
		"label":    "【アンケート】商談予定額　（TOPIC）",
		"label_en": "[Survey] Planned Business Negotiation Amount　（TOPIC）",
		"type":     "number",
		"options":  []map[string]any{},
	},
}

var TOPIC_VISITOR_INDUSTRY_OPTIONS_FOODEX = []map[string]any{
	map[string]any{
		"value":    1,
		"label":    "商社・卸 ➞ 卸・問屋",
		"label_en": "Trading / Wholesale ➞ Food Trading Company / Wholesaler",
	},
	map[string]any{
		"value":    2,
		"label":    "商社・卸 ➞ 輸入商社",
		"label_en": "Trading / Wholesale ➞ Importer: Foreign food products to Japan",
	},
	map[string]any{
		"value":    43,
		"label":    "商社・卸 ➞ 輸出商社",
		"label_en": "Trading / Wholesale ➞ Exporter: Japanese food products to overseas",
	},
	map[string]any{
		"value":    3,
		"label":    "商社・卸 ➞ その他の中間流通",
		"label_en": "Trading / Wholesale ➞ Intermediary Distributor",
	},
	map[string]any{
		"value":    4,
		"label":    "外食・給食・中食 ➞ 外食",
		"label_en": "Food Service ➞ Restaurants, Foodservice",
	},
	map[string]any{
		"value":    5,
		"label":    "外食・給食・中食 ➞ 給食",
		"label_en": "Food Service ➞ Lunchroom / Catering",
	},
	map[string]any{
		"value":    6,
		"label":    "外食・給食・中食 ➞ 中食",
		"label_en": "Food Service ➞ Delicatessen",
	},
	map[string]any{
		"value":    7,
		"label":    "外食・給食・中食 ➞ その他のフードサービス",
		"label_en": "Food Service ➞ Other Foodservices",
	},
	map[string]any{
		"value":    8,
		"label":    "小売 ➞ 百貨店",
		"label_en": "Retailer ➞ Department Store",
	},
	map[string]any{
		"value":    9,
		"label":    "小売 ➞ スーパーマーケット",
		"label_en": "Retailer ➞ Supermarket",
	},
	map[string]any{
		"value":    10,
		"label":    "小売 ➞ コンビニエンスストア",
		"label_en": "Retailer ➞ Convenience Store",
	},
	map[string]any{
		"value":    11,
		"label":    "小売 ➞ ディスカウントストア",
		"label_en": "Retailer ➞ Discount Store",
	},
	map[string]any{
		"value":    12,
		"label":    "小売 ➞ ドラッグストア",
		"label_en": "Retailer ➞ Pharmacy",
	},
	map[string]any{
		"value":    13,
		"label":    "小売 ➞ 専門店",
		"label_en": "Retailer ➞ Specialty Shop",
	},
	map[string]any{
		"value":    14,
		"label":    "小売 ➞ 通販・ネットスーパー",
		"label_en": "Retailer ➞ E-Commerce",
	},
	map[string]any{
		"value":    15,
		"label":    "ホテル・旅館・レジャー ➞ ホテル・旅館",
		"label_en": "Hotels / Inns / Leisure Facilities ➞ Hotel / Inn",
	},
	map[string]any{
		"value":    16,
		"label":    "ホテル・旅館・レジャー ➞ 式場・会館",
		"label_en": "Hotels / Inns / Leisure Facilities ➞ Banquet, Wedding Facility / Hall",
	},
	map[string]any{
		"value":    17,
		"label":    "ホテル・旅館・レジャー ➞ レジャー",
		"label_en": "Hotels / Inns / Leisure Facilities ➞ Leisure Facility",
	},
	map[string]any{
		"value":    18,
		"label":    "食品・飲料メーカー／生産者 ➞ 食品/飲料メーカー",
		"label_en": "Food / Beverage Manufacture ➞ Food and Beverage Manufacturers",
	},
	map[string]any{
		"value":    19,
		"label":    "食品・飲料メーカー／生産者 ➞ 食品受託加工メーカー",
		"label_en": "Food / Beverage Manufacture ➞ OEM Manufacturers",
	},
	map[string]any{
		"value":    20,
		"label":    "食品・飲料メーカー／生産者 ➞ その他の食品メーカー",
		"label_en": "Food / Beverage Manufacture ➞ Other Food and Beverage Manufacturers",
	},
	map[string]any{
		"value":    21,
		"label":    "食品・飲料メーカー／生産者 ➞ 生産者(農業・畜産・水産業)",
		"label_en": "Food / Beverage Manufacture ➞ Producers (Agriculture, Livestock, Fisheries)",
	},
	map[string]any{
		"value":    23,
		"label":    "機器・設備メーカー ➞ 食品加工機器メーカー",
		"label_en": "Manufacturer ➞ Food Processing Equipment Manufacturer",
	},
	map[string]any{
		"value":    24,
		"label":    "機器・設備メーカー ➞ 厨房・調理機器メーカー",
		"label_en": "Manufacturer ➞ Cooking Utensil Manufacturer",
	},
	map[string]any{
		"value":    25,
		"label":    "機器・設備メーカー ➞ 衛生・消耗品メーカー",
		"label_en": "Manufacturer ➞ Sanitation / Cleaning Equipment Manufacturer",
	},
	map[string]any{
		"value":    26,
		"label":    "機器・設備メーカー ➞ 包装・容器メーカー",
		"label_en": "Manufacturer ➞ Packaging / Container Manufacturer",
	},
	map[string]any{
		"value":    27,
		"label":    "機器・設備メーカー ➞ その他の機器・設備メーカー",
		"label_en": "Manufacturer ➞ Other Equipment Manufacture",
	},
	map[string]any{
		"value":    28,
		"label":    "官公庁・団体・専門家 ➞ 官公庁・政府関係/公的機関",
		"label_en": "Government Body / Association / Specialist ➞ Government Body",
	},
	map[string]any{
		"value":    45,
		"label":    "官公庁・団体・専門家 ➞ その他団体",
		"label_en": "Government Body / Association / Specialist ➞ Association",
	},
	map[string]any{
		"value":    46,
		"label":    "官公庁・団体・専門家 ➞ 専門家",
		"label_en": "Government Body / Association / Specialist ➞ Specialist",
	},
	map[string]any{
		"value":    29,
		"label":    "官公庁・団体・専門家 ➞ 教育・教育機関",
		"label_en": "Government Body / Association / Specialist ➞ Education / Educational Institutional",
	},
	map[string]any{
		"value":    31,
		"label":    "物流・運輸 ➞ 物流業",
		"label_en": "Logistics ➞ Logistics",
	},
	map[string]any{
		"value":    32,
		"label":    "物流・運輸 ➞ 倉庫業",
		"label_en": "Logistics ➞ Warehouse",
	},
	map[string]any{
		"value":    33,
		"label":    "エステ・スパ ➞ エステ・スパ",
		"label_en": "Beauty Salon / Spa ➞ Beauty Salon / Spa",
	},
	map[string]any{
		"value":    34,
		"label":    "印刷会社 ➞ 印刷会社",
		"label_en": "Printing Company ➞ Printing Company",
	},
	map[string]any{
		"value":    35,
		"label":    "コンサルティング ➞ コンサルティング",
		"label_en": "Consultancy ➞ Consultancy",
	},
	map[string]any{
		"value":    36,
		"label":    "金融業・保険業 ➞ 金融業・保険業",
		"label_en": "Financial / Insurance Business ➞ Financial / Insurance Business",
	},
	map[string]any{
		"value":    37,
		"label":    "旅行業・旅行代理店 ➞ 旅行業・旅行代理店",
		"label_en": "Travel Agency ➞ Travel Agency",
	},
	map[string]any{
		"value":    38,
		"label":    "アドバイザー ➞ アドバイザー",
		"label_en": "Adviser ➞ Adviser",
	},
	map[string]any{
		"value":    39,
		"label":    "主婦・主夫 ➞ 主婦・主夫",
		"label_en": "Homemaker ➞ Homemaker",
	},
	map[string]any{
		"value":    40,
		"label":    "出版 ➞ 出版",
		"label_en": "Publishing Company ➞ Publishing Company",
	},
	map[string]any{
		"value":    41,
		"label":    "報道関係者 ➞ 報道関係者",
		"label_en": "Press ➞ Press",
	},
	map[string]any{
		"value":    44,
		"label":    "料理専門家 ➞ 料理専門家",
		"label_en": "Culinary Professionals ➞ Culinary Professionals",
	},
	map[string]any{"value": 42, "label": "その他 ➞ その他", "label_en": "Other ➞ Other"},
}

var TOPIC_VISITOR_INDUSTRY_OPTIONS_HCJ = []map[string]any{
	map[string]any{
		"value":    1,
		"label":    "宿泊 ➞ シティホテル",
		"label_en": "Hotels / Inns ➞ City Hotel",
	},
	map[string]any{
		"value":    2,
		"label":    "宿泊 ➞ ビジネスホテル",
		"label_en": "Hotels / Inns ➞ Business Hotel",
	},
	map[string]any{
		"value":    3,
		"label":    "宿泊 ➞ リゾートホテル",
		"label_en": "Hotels / Inns ➞ Resort Hotel",
	},
	map[string]any{
		"value":    4,
		"label":    "宿泊 ➞ 外資系ホテル",
		"label_en": "Hotels / Inns ➞ Foreign Capital Hotel",
	},
	map[string]any{
		"value":    5,
		"label":    "宿泊 ➞ その他のホテル",
		"label_en": "Hotels / Inns ➞ Other Hotel",
	},
	map[string]any{
		"value":    6,
		"label":    "宿泊 ➞ 旅館",
		"label_en": "Hotels / Inns ➞ Ryokan / Japanese Style Hotel",
	},
	map[string]any{
		"value":    7,
		"label":    "宿泊 ➞ ペンション / 民宿",
		"label_en": "Hotels / Inns ➞ Inn / B&B / Guest House",
	},
	map[string]any{
		"value":    8,
		"label":    "宿泊 ➞ グランピング / キャンプ / コテージ",
		"label_en": "Hotels / Inns ➞ Glamping / Camping / Cottage",
	},
	map[string]any{"value": 9, "label": "宿泊 ➞ その他", "label_en": "Hotels / Inns ➞ Others"},
	map[string]any{
		"value":    10,
		"label":    "ウェディング・スポーツ・レジャー ➞ ウェディング・会館・ホール",
		"label_en": "Wedding / Sports / Leisure ➞ Wedding / Banquet / Hall",
	},
	map[string]any{
		"value":    11,
		"label":    "ウェディング・スポーツ・レジャー ➞ フィットネス / スポーツ施設",
		"label_en": "Wedding / Sports / Leisure ➞ Fitness / Sports Facility",
	},
	map[string]any{
		"value":    12,
		"label":    "ウェディング・スポーツ・レジャー ➞ レジャー施設",
		"label_en": "Wedding / Sports / Leisure ➞ Leisure",
	},
	map[string]any{
		"value":    13,
		"label":    "ウェディング・スポーツ・レジャー ➞ その他",
		"label_en": "Wedding / Sports / Leisure ➞ Others",
	},
	map[string]any{
		"value":    14,
		"label":    "外食 ➞ ファミリーレストラン",
		"label_en": "Restaurants / Food Services ➞ Family Restaurant / Diner",
	},
	map[string]any{
		"value":    15,
		"label":    "外食 ➞ 日本料理・割烹・料亭",
		"label_en": "Restaurants / Food Services ➞ Traditional Japanese Restaurant",
	},
	map[string]any{
		"value":    16,
		"label":    "外食 ➞ 寿司",
		"label_en": "Restaurants / Food Services ➞ Sushi",
	},
	map[string]any{
		"value":    17,
		"label":    "外食 ➞ そば・うどん",
		"label_en": "Restaurants / Food Services ➞ Japanese Noodle",
	},
	map[string]any{
		"value":    18,
		"label":    "外食 ➞ らーめん",
		"label_en": "Restaurants / Food Services ➞ Ramen Noodle",
	},
	map[string]any{
		"value":    19,
		"label":    "外食 ➞ 焼肉",
		"label_en": "Restaurants / Food Services ➞ Barbecue",
	},
	map[string]any{
		"value":    20,
		"label":    "外食 ➞ 中華料理",
		"label_en": "Restaurants / Food Services ➞ Chinese",
	},
	map[string]any{
		"value":    21,
		"label":    "外食 ➞ イタリア料理",
		"label_en": "Restaurants / Food Services ➞ Italian",
	},
	map[string]any{
		"value":    22,
		"label":    "外食 ➞ フランス料理",
		"label_en": "Restaurants / Food Services ➞ French",
	},
	map[string]any{
		"value":    23,
		"label":    "外食 ➞ アジア・エスニック料理",
		"label_en": "Restaurants / Food Services ➞ Asian / Ethnic",
	},
	map[string]any{
		"value":    24,
		"label":    "外食 ➞ カフェ・喫茶店",
		"label_en": "Restaurants / Food Services ➞ Cafe / Tearoom",
	},
	map[string]any{
		"value":    25,
		"label":    "外食 ➞ 居酒屋",
		"label_en": "Restaurants / Food Services ➞ Japanese Bar Restaurant",
	},
	map[string]any{
		"value":    26,
		"label":    "外食 ➞ Bar・パブ",
		"label_en": "Restaurants / Food Services ➞ Bar / Pub",
	},
	map[string]any{
		"value":    27,
		"label":    "外食 ➞ ファストフード",
		"label_en": "Restaurants / Food Services ➞ Fast Food",
	},
	map[string]any{
		"value":    28,
		"label":    "外食 ➞ ベーカリー・菓子",
		"label_en": "Restaurants / Food Services ➞ Bakery / Confectionery",
	},
	map[string]any{
		"value":    29,
		"label":    "外食 ➞ その他",
		"label_en": "Restaurants / Food Services ➞ Others",
	},
	map[string]any{
		"value":    30,
		"label":    "給食 ➞ 産業給食",
		"label_en": "Meal Services ➞ Industrial Catering",
	},
	map[string]any{
		"value":    31,
		"label":    "給食 ➞ 病院給食",
		"label_en": "Meal Services ➞ Hospital Food Service",
	},
	map[string]any{
		"value":    32,
		"label":    "給食 ➞ 学校給食",
		"label_en": "Meal Services ➞ School Food Service",
	},
	map[string]any{
		"value":    33,
		"label":    "給食 ➞ 福祉給食",
		"label_en": "Meal Services ➞ Welfare Food Service",
	},
	map[string]any{"value": 34, "label": "給食 ➞ その他", "label_en": "Meal Services ➞ Others"},
	map[string]any{
		"value":    35,
		"label":    "中食 ➞ 惣菜",
		"label_en": "Delicatessen / Catering Services ➞ Delicatessen",
	},
	map[string]any{
		"value":    36,
		"label":    "中食 ➞ 弁当",
		"label_en": "Delicatessen / Catering Services ➞ Bento / Box lunch",
	},
	map[string]any{
		"value":    37,
		"label":    "中食 ➞ 宅配",
		"label_en": "Delicatessen / Catering Services ➞ Meal Delivery Service",
	},
	map[string]any{
		"value":    38,
		"label":    "中食 ➞ その他",
		"label_en": "Delicatessen / Catering Services ➞ Others",
	},
	map[string]any{"value": 39, "label": "小売 ➞ 百貨店", "label_en": "Retailers ➞ Department Store"},
	map[string]any{
		"value":    40,
		"label":    "小売 ➞ スーパーマーケット",
		"label_en": "Retailers ➞ Supermarket",
	},
	map[string]any{
		"value":    41,
		"label":    "小売 ➞ コンビニエンスストア",
		"label_en": "Retailers ➞ Convenience Store",
	},
	map[string]any{
		"value":    42,
		"label":    "小売 ➞ ディスカウントストア",
		"label_en": "Retailers ➞ Discount Store",
	},
	map[string]any{
		"value":    43,
		"label":    "小売 ➞ ドラッグストア",
		"label_en": "Retailers ➞ Drug Store",
	},
	map[string]any{"value": 44, "label": "小売 ➞ 専門店", "label_en": "Retailers ➞ Specialty Shop"},
	map[string]any{
		"value":    45,
		"label":    "小売 ➞ 通販・ネットスーパー",
		"label_en": "Retailers ➞ EC / Online Grocery Store",
	},
	map[string]any{"value": 46, "label": "小売 ➞ その他", "label_en": "Retailers ➞ Others"},
	map[string]any{
		"value":    47,
		"label":    "スパ・温浴・サウナ ➞ スパ・温泉・温浴施設",
		"label_en": "Spa / Bath / Sauna ➞ Spa / Hot Spring / Hot Bath",
	},
	map[string]any{
		"value":    48,
		"label":    "スパ・温浴・サウナ ➞ エステ・リラクゼーション",
		"label_en": "Spa / Bath / Sauna ➞ Esthetic / Relaxation Salon",
	},
	map[string]any{
		"value":    49,
		"label":    "スパ・温浴・サウナ ➞ その他",
		"label_en": "Spa / Bath / Sauna ➞ Others",
	},
	map[string]any{
		"value":    50,
		"label":    "商社・卸・問屋 ➞ 厨房機器・調理器具",
		"label_en": "Tradings / Wholesales ➞ Food Service Equipment",
	},
	map[string]any{
		"value":    51,
		"label":    "商社・卸・問屋 ➞ 食品・飲料・調味料",
		"label_en": "Tradings / Wholesales ➞ Food / Beverage / Seasoning",
	},
	map[string]any{
		"value":    52,
		"label":    "商社・卸・問屋 ➞ 家具・インテリア・エクステリア",
		"label_en": "Tradings / Wholesales ➞ Furniture / Interior / Exterior",
	},
	map[string]any{
		"value":    53,
		"label":    "商社・卸・問屋 ➞ テーブルウエア・食品容器",
		"label_en": "Tradings / Wholesales ➞ Tableware / Food Container",
	},
	map[string]any{
		"value":    54,
		"label":    "商社・卸・問屋 ➞ 温浴・サウナ関連",
		"label_en": "Tradings / Wholesales ➞ Spa / Bath / Sauna",
	},
	map[string]any{
		"value":    55,
		"label":    "商社・卸・問屋 ➞ テクノロジー(IoT/ICT/AI/ロボット)",
		"label_en": "Tradings / Wholesales ➞ Technology (IoT / ICT / AI / Robot)",
	},
	map[string]any{
		"value":    56,
		"label":    "商社・卸・問屋 ➞ アメニティ・消耗品",
		"label_en": "Tradings / Wholesales ➞ Amenity / Consumable",
	},
	map[string]any{
		"value":    57,
		"label":    "商社・卸・問屋 ➞ 衛生・清掃品",
		"label_en": "Tradings / Wholesales ➞ Sanitation / Cleaning Equipment",
	},
	map[string]any{
		"value":    58,
		"label":    "商社・卸・問屋 ➞ 電気機械・器具",
		"label_en": "Tradings / Wholesales ➞ Electric Appliance",
	},
	map[string]any{
		"value":    59,
		"label":    "商社・卸・問屋 ➞ その他",
		"label_en": "Tradings / Wholesales ➞ Others",
	},
	map[string]any{
		"value":    60,
		"label":    "開発・設計・建築・デザイン ➞ 設計・デザイン",
		"label_en": "Estate Development / Construction / Design ➞ Construction / Design",
	},
	map[string]any{
		"value":    61,
		"label":    "開発・設計・建築・デザイン ➞ 空間コーディネーター",
		"label_en": "Estate Development / Construction / Design ➞ Space Coordinator",
	},
	map[string]any{
		"value":    62,
		"label":    "開発・設計・建築・デザイン ➞ 建設・建築・工務店",
		"label_en": "Estate Development / Construction / Design ➞ Construction / Architecture",
	},
	map[string]any{
		"value":    63,
		"label":    "開発・設計・建築・デザイン ➞ 不動産ディベロッパー",
		"label_en": "Estate Development / Construction / Design ➞ Estate Development",
	},
	map[string]any{
		"value":    64,
		"label":    "開発・設計・建築・デザイン ➞ その他",
		"label_en": "Estate Development / Construction / Design ➞ Others",
	},
	map[string]any{
		"value":    65,
		"label":    "メーカー ➞ 厨房機器・調理器具",
		"label_en": "Manufacturers ➞ Food Service Equipment",
	},
	map[string]any{
		"value":    66,
		"label":    "メーカー ➞ 食品・飲料・調味料",
		"label_en": "Manufacturers ➞ Food / Beverage / Seasoning",
	},
	map[string]any{
		"value":    67,
		"label":    "メーカー ➞ 家具・インテリア・エクステリア",
		"label_en": "Manufacturers ➞ Furniture / Interior / Exterior",
	},
	map[string]any{
		"value":    68,
		"label":    "メーカー ➞ テーブルウエア・食品容器（使い捨て以外）",
		"label_en": "Manufacturers ➞ Tableware / Food Container (non-Disposable)",
	},
	map[string]any{
		"value":    69,
		"label":    "メーカー ➞ 包装・容器（使い捨て）",
		"label_en": "Manufacturers ➞ Package (Disposable)",
	},
	map[string]any{
		"value":    70,
		"label":    "メーカー ➞ 温浴・サウナ関連",
		"label_en": "Manufacturers ➞ Spa / Bath / Sauna",
	},
	map[string]any{
		"value":    71,
		"label":    "メーカー ➞ テクノロジー(IoT/ICT/AI/ロボット)",
		"label_en": "Manufacturers ➞ Technology (IoT / ICT / AI / Robot)",
	},
	map[string]any{
		"value":    72,
		"label":    "メーカー ➞ アメニティ・消耗品",
		"label_en": "Manufacturers ➞ Amenity / Consumable",
	},
	map[string]any{
		"value":    73,
		"label":    "メーカー ➞ 衛生・清掃品",
		"label_en": "Manufacturers ➞ Sanitation / Cleaning Equipment",
	},
	map[string]any{
		"value":    74,
		"label":    "メーカー ➞ 電気機械・器具",
		"label_en": "Manufacturers ➞ Electric Appliance",
	},
	map[string]any{"value": 75, "label": "メーカー ➞ 繊維", "label_en": "Manufacturers ➞ Textile"},
	map[string]any{
		"value":    76,
		"label":    "メーカー ➞ 通信",
		"label_en": "Manufacturers ➞ Communication",
	},
	map[string]any{"value": 77, "label": "メーカー ➞ その他", "label_en": "Manufacturers ➞ Others"},
	map[string]any{
		"value":    78,
		"label":    "官公庁・自治体・団体 ➞ 官公庁",
		"label_en": "Government / Local Government / Organization ➞ Government",
	},
	map[string]any{
		"value":    79,
		"label":    "官公庁・自治体・団体 ➞ 自治体",
		"label_en": "Government / Local Government / Organization ➞ Local Government",
	},
	map[string]any{
		"value":    80,
		"label":    "官公庁・自治体・団体 ➞ 団体",
		"label_en": "Government / Local Government / Organization ➞ Organization",
	},
	map[string]any{
		"value":    81,
		"label":    "官公庁・自治体・団体 ➞ 教育・研究機関",
		"label_en": "Government / Local Government / Organization ➞ Education / Research Institute",
	},
	map[string]any{
		"value":    82,
		"label":    "官公庁・自治体・団体 ➞ その他",
		"label_en": "Government / Local Government / Organization ➞ Others",
	},
	map[string]any{
		"value":    83,
		"label":    "ビルメンテナンス ➞ 総合管理",
		"label_en": "Building Maintenances ➞ Comprehensive Management",
	},
	map[string]any{
		"value":    84,
		"label":    "ビルメンテナンス ➞ 清掃・衛生管理",
		"label_en": "Building Maintenances ➞ Cleaning / Hygiene Management",
	},
	map[string]any{
		"value":    85,
		"label":    "ビルメンテナンス ➞ 設備管理",
		"label_en": "Building Maintenances ➞ Equipment Management",
	},
	map[string]any{
		"value":    86,
		"label":    "ビルメンテナンス ➞ 警備・セキュリティ",
		"label_en": "Building Maintenances ➞ Security",
	},
	map[string]any{
		"value":    87,
		"label":    "ビルメンテナンス ➞ その他",
		"label_en": "Building Maintenances ➞ Others",
	},
	map[string]any{
		"value":    88,
		"label":    "その他関連業種 ➞ 農業・林業・漁業",
		"label_en": "Other Industries ➞ Agriculture / Forestry / Fishery",
	},
	map[string]any{
		"value":    89,
		"label":    "その他関連業種 ➞ 病院・福祉施設",
		"label_en": "Other Industries ➞ Hospital / Welfare Facility",
	},
	map[string]any{
		"value":    90,
		"label":    "その他関連業種 ➞ 旅行代理店",
		"label_en": "Other Industries ➞ Travel Agency",
	},
	map[string]any{
		"value":    91,
		"label":    "その他関連業種 ➞ 空港・運輸・運送",
		"label_en": "Other Industries ➞ Airport / Transportation / Logistics",
	},
	map[string]any{
		"value":    92,
		"label":    "その他関連業種 ➞ その他",
		"label_en": "Other Industries ➞ Others",
	},
	map[string]any{"value": 93, "label": "学生 ➞ 学生", "label_en": "Students ➞ Students"},
}

var TOPIC_VISITOR_JOB_OPTIONS_FOODEX = []map[string]any{
	map[string]any{"value": 1, "label": "経営", "label_en": "Management"},
	map[string]any{"value": 2, "label": "営業・販売", "label_en": "Sales & Distribution"},
	map[string]any{"value": 3, "label": "仕入・購買", "label_en": "Purchasing"},
	map[string]any{
		"value":    4,
		"label":    "商品企画・マーケティング",
		"label_en": "Planning / Marketing",
	},
	map[string]any{"value": 5, "label": "貿易", "label_en": "Overseas Trade"},
	map[string]any{"value": 6, "label": "研究開発・品質管理", "label_en": "R&D / Quality Control"},
	map[string]any{
		"value":    7,
		"label":    "総務・人事・広報",
		"label_en": "General Affairs / Human Resources / PR",
	},
	map[string]any{
		"value":    8,
		"label":    "調理・加工・サービス",
		"label_en": "Cooking, Processing & Services",
	},
	map[string]any{"value": 99, "label": "その他", "label_en": "Other"},
}

var TOPIC_VISITOR_JOB_OPTIONS_HCJ = []map[string]any{
	map[string]any{"value": 1, "label": "経営・役員", "label_en": "CEO / President / Director"},
	map[string]any{"value": 2, "label": "仕入・購買", "label_en": "Purchasing"},
	map[string]any{"value": 3, "label": "調理・栄養士", "label_en": "Cooks / Nutritionists"},
	map[string]any{"value": 4, "label": "営業", "label_en": "Sales"},
	map[string]any{"value": 5, "label": "マーケティング", "label_en": "Marketing"},
	map[string]any{
		"value":    6,
		"label":    "総務・人事・経営企画",
		"label_en": "General Affairs / HR / Corporate Planning",
	},
	map[string]any{
		"value":    7,
		"label":    "開発・店舗開発",
		"label_en": "Development / Store Development",
	},
	map[string]any{"value": 8, "label": "システム", "label_en": "Systems"},
	map[string]any{
		"value":    9,
		"label":    "製造・生産技術",
		"label_en": "Manufacturing / Production Technology",
	},
	map[string]any{"value": 10, "label": "広報・宣伝", "label_en": "Public Relations / Advertising"},
	map[string]any{
		"value":    11,
		"label":    "品質管理・施設管理",
		"label_en": "Quality Control / Facility Management",
	},
	map[string]any{"value": 12, "label": "接客・サービス", "label_en": "Customer Service"},
	map[string]any{"value": 13, "label": "設計・デザイン", "label_en": "Construction / Design"},
	map[string]any{"value": 14, "label": "その他", "label_en": "Others"},
}

var TOPIC_VISITOR_POSITION_OPTIONS_FOODEX = []map[string]any{
	map[string]any{"value": 1, "label": "経営・役員", "label_en": "CEO / President / Director"},
	map[string]any{"value": 2, "label": "部長・室長", "label_en": "Department / Division Director"},
	map[string]any{"value": 3, "label": "次長・課長", "label_en": "Section Head / Chief"},
	map[string]any{"value": 4, "label": "係長・主任", "label_en": "Manager"},
	map[string]any{"value": 5, "label": "専門職", "label_en": "Specialist"},
	map[string]any{"value": 6, "label": "一般社員", "label_en": "Employee"},
}

var TOPIC_VISITOR_POSITION_OPTIONS_HCJ = []map[string]any{
	map[string]any{"value": 1, "label": "経営・役員", "label_en": "CEO / President / Director"},
	map[string]any{"value": 2, "label": "部長・室長", "label_en": "General Manager / Chief Manager"},
	map[string]any{"value": 3, "label": "次長・課長", "label_en": "Section Chief / Manager"},
	map[string]any{"value": 4, "label": "係長・主任", "label_en": "Unit Head / Deputy Manager"},
	map[string]any{"value": 5, "label": "専門職", "label_en": "Specialist"},
	map[string]any{"value": 6, "label": "一般社員", "label_en": "Staff / Employee"},
}

var TOPIC_VISITOR_INTERESTED_CATEGORIES_FOODEX = []map[string]any{
	map[string]any{
		"value":    1,
		"label":    "ベーカリー, ケーキ＆デザート ➞ 焼き菓子",
		"label_en": "Bakery, Cakes & Desserts ➞ Baked Goods",
	},
	map[string]any{
		"value":    2,
		"label":    "ベーカリー, ケーキ＆デザート ➞ ケーキ / デザート",
		"label_en": "Bakery, Cakes & Desserts ➞ Cakes / Dessert",
	},
	map[string]any{
		"value":    3,
		"label":    "ベーカリー, ケーキ＆デザート ➞ パン",
		"label_en": "Bakery, Cakes & Desserts ➞ Bread",
	},
	map[string]any{
		"value":    4,
		"label":    "ベーカリー, ケーキ＆デザート ➞ 和菓子",
		"label_en": "Bakery, Cakes & Desserts ➞ Japanese sweets",
	},
	map[string]any{
		"value":    5,
		"label":    "ベーカリー, ケーキ＆デザート ➞ グルテンフリー",
		"label_en": "Bakery, Cakes & Desserts ➞ Gluten-free",
	},
	map[string]any{
		"value":    129,
		"label":    "ベーカリー, ケーキ＆デザート ➞ チョコレート",
		"label_en": "Bakery, Cakes & Desserts ➞ Chocolate",
	},
	map[string]any{
		"value":    130,
		"label":    "ベーカリー, ケーキ＆デザート ➞ 飴",
		"label_en": "Bakery, Cakes & Desserts ➞ Candy",
	},
	map[string]any{
		"value":    6,
		"label":    "飲料 ➞ ボトルウォーター",
		"label_en": "Beverages ➞ Bottled Water",
	},
	map[string]any{
		"value":    7,
		"label":    "飲料 ➞ 炭酸飲料",
		"label_en": "Beverages ➞ Carbonated Soft Drinks",
	},
	map[string]any{"value": 8, "label": "飲料 ➞ コーヒー", "label_en": "Beverages ➞ Coffee"},
	map[string]any{
		"value":    9,
		"label":    "飲料 ➞ エナジードリンク",
		"label_en": "Beverages ➞ Energy Drinks",
	},
	map[string]any{
		"value":    10,
		"label":    "飲料 ➞ フレーバーウォーター",
		"label_en": "Beverages ➞ Flavoured Water",
	},
	map[string]any{
		"value":    11,
		"label":    "飲料 ➞ インスタント飲料",
		"label_en": "Beverages ➞ Instant Beverages",
	},
	map[string]any{"value": 12, "label": "飲料 ➞ 果汁飲料", "label_en": "Beverages ➞ Juice"},
	map[string]any{
		"value":    13,
		"label":    "飲料 ➞ ノンアルコール飲料",
		"label_en": "Beverages ➞ Non-alcoholic Beverages",
	},
	map[string]any{"value": 14, "label": "飲料 ➞ スムージー", "label_en": "Beverages ➞ Smoothies"},
	map[string]any{"value": 15, "label": "飲料 ➞ 茶", "label_en": "Beverages ➞ Tea"},
	map[string]any{
		"value":    16,
		"label":    "生鮮食品 ➞ 卵・卵加工品",
		"label_en": "Fresh Food ➞ Eggs & Egg Products",
	},
	map[string]any{"value": 17, "label": "生鮮食品 ➞ キノコ類", "label_en": "Fresh Food ➞ Mushroom"},
	map[string]any{
		"value":    18,
		"label":    "生鮮食品 ➞ 果物・野菜",
		"label_en": "Fresh Food ➞ Fresh Fruit & Vegetables",
	},
	map[string]any{
		"value":    124,
		"label":    "生鮮食品 ➞ 農産加工品",
		"label_en": "Fresh Food ➞ Processed agricultural products",
	},
	map[string]any{
		"value":    19,
		"label":    "調味料＆ソース ➞ 醤油",
		"label_en": "Condiments & Sauces ➞ Soy Sauce",
	},
	map[string]any{
		"value":    20,
		"label":    "調味料＆ソース ➞ 味噌",
		"label_en": "Condiments & Sauces ➞ Miso",
	},
	map[string]any{
		"value":    21,
		"label":    "調味料＆ソース ➞ 粉末調味料",
		"label_en": "Condiments & Sauces ➞ Powdered Seasoning",
	},
	map[string]any{
		"value":    22,
		"label":    "調味料＆ソース ➞ 液体調味料",
		"label_en": "Condiments & Sauces ➞ Liquid Seasonings",
	},
	map[string]any{
		"value":    104,
		"label":    "調味料＆ソース ➞ 固形調味料",
		"label_en": "Condiments & Sauces ➞ Solid Seasonings",
	},
	map[string]any{
		"value":    105,
		"label":    "調味料＆ソース ➞ ふりかけ、茶漬など",
		"label_en": "Condiments & Sauces ➞ Rice Toppings & Ochazuke Mixes",
	},
	map[string]any{
		"value":    23,
		"label":    "調味料＆ソース ➞ 酢",
		"label_en": "Condiments & Sauces ➞ Vinegar",
	},
	map[string]any{
		"value":    24,
		"label":    "乾物＆インスタント・レトルト食品 ➞ 缶詰",
		"label_en": "Dry Foods & Instant Foods ➞ Canned Products",
	},
	map[string]any{
		"value":    25,
		"label":    "乾物＆インスタント・レトルト食品 ➞ ドライフルーツ",
		"label_en": "Dry Foods & Instant Foods ➞ Dried Fruits",
	},
	map[string]any{
		"value":    106,
		"label":    "乾物＆インスタント・レトルト食品 ➞ レトルトパウチ食品、密封包装食品（カレー、丼の具など）",
		"label_en": "Dry Foods & Instant Foods ➞ Ready-to-eat / Sealed Packaged Foods",
	},
	map[string]any{
		"value":    28,
		"label":    "乾物＆インスタント・レトルト食品 ➞ 乾燥水産品",
		"label_en": "Dry Foods & Instant Foods ➞ Dry Seafood",
	},
	map[string]any{
		"value":    29,
		"label":    "乾物＆インスタント・レトルト食品 ➞ パスタ",
		"label_en": "Dry Foods & Instant Foods ➞ Pasta",
	},
	map[string]any{
		"value":    30,
		"label":    "冷蔵食品・乳製品 ➞ 代替乳製品",
		"label_en": "Chilled Food/Daily ➞ Alternative Milk Products",
	},
	map[string]any{
		"value":    31,
		"label":    "冷蔵食品・乳製品 ➞ バター",
		"label_en": "Chilled Food/Daily ➞ Butter",
	},
	map[string]any{
		"value":    32,
		"label":    "冷蔵食品・乳製品 ➞ チーズ",
		"label_en": "Chilled Food/Daily ➞ Cheese",
	},
	map[string]any{
		"value":    33,
		"label":    "冷蔵食品・乳製品 ➞ クリーム",
		"label_en": "Chilled Food/Daily ➞ Cream & Cream Products",
	},
	map[string]any{
		"value":    34,
		"label":    "冷蔵食品・乳製品 ➞ 牛乳",
		"label_en": "Chilled Food/Daily ➞ Milk",
	},
	map[string]any{
		"value":    35,
		"label":    "冷蔵食品・乳製品 ➞ ヨーグルト（レギュラー、フレーバー、フルーツなど）",
		"label_en": "Chilled Food/Daily ➞ Yoghurt (Regular, Flavoured, Fruit, etc.)",
	},
	map[string]any{
		"value":    36,
		"label":    "冷蔵食品・乳製品 ➞ 豆腐・納豆",
		"label_en": "Chilled Food/Daily ➞ Tofu/Natto",
	},
	map[string]any{
		"value":    37,
		"label":    "機能性食品・健康食品・代替食 ➞ 機能性食品",
		"label_en": "Functional Foods, Health Foods, Alternative Foods ➞ Fortified & Functional Products",
	},
	map[string]any{
		"value":    38,
		"label":    "機能性食品・健康食品・代替食 ➞ プロテイン飲料",
		"label_en": "Functional Foods, Health Foods, Alternative Foods ➞ Protein Drink",
	},
	map[string]any{
		"value":    39,
		"label":    "機能性食品・健康食品・代替食 ➞ サプリメント",
		"label_en": "Functional Foods, Health Foods, Alternative Foods ➞ Supplements",
	},
	map[string]any{
		"value":    40,
		"label":    "機能性食品・健康食品・代替食 ➞ ヘルス＆ウェルネス製品",
		"label_en": "Functional Foods, Health Foods, Alternative Foods ➞ Health & Wellness Products",
	},
	map[string]any{
		"value":    108,
		"label":    "機能性食品・健康食品・代替食 ➞ プラントベース／代替食",
		"label_en": "Functional Foods, Health Foods, Alternative Foods ➞ Plant-based / Alternative Foods",
	},
	map[string]any{
		"value":    125,
		"label":    "機能性食品・健康食品・代替食 ➞ 健康飲料",
		"label_en": "Functional Foods, Health Foods, Alternative Foods ➞ Health drinks",
	},
	map[string]any{
		"value":    41,
		"label":    "冷凍食品 ➞ 冷凍焼き菓子",
		"label_en": "Frozen Food ➞ Frozen Baked Goods",
	},
	map[string]any{
		"value":    42,
		"label":    "冷凍食品 ➞ 冷凍乳製品",
		"label_en": "Frozen Food ➞ Frozen Dairy Products",
	},
	map[string]any{
		"value":    43,
		"label":    "冷凍食品 ➞ 冷凍野菜・果物",
		"label_en": "Frozen Food ➞ Frozen Fruit & Vegetables",
	},
	map[string]any{
		"value":    44,
		"label":    "冷凍食品 ➞ 冷凍肉",
		"label_en": "Frozen Food ➞ Frozen Meat",
	},
	map[string]any{
		"value":    45,
		"label":    "冷凍食品 ➞ 冷凍水産品",
		"label_en": "Frozen Food ➞ Frozen Seafood",
	},
	map[string]any{
		"value":    46,
		"label":    "冷凍食品 ➞ 冷凍調理済み食品",
		"label_en": "Frozen Food ➞ Frozen Ready Meals",
	},
	map[string]any{
		"value":    47,
		"label":    "冷凍食品 ➞ アイスクリーム／シャーベット・ジェラート",
		"label_en": "Frozen Food ➞ Ice Cream / Sorbet Gelato",
	},
	map[string]any{
		"value":    48,
		"label":    "ベジタリアン・ヴィーガン・ハラール・コーシャ ➞ ベジタリアン",
		"label_en": "Vegitarian, Vegan, Halal and Kosher Food ➞ Vegetarian",
	},
	map[string]any{
		"value":    49,
		"label":    "ベジタリアン・ヴィーガン・ハラール・コーシャ ➞ ヴィーガン",
		"label_en": "Vegitarian, Vegan, Halal and Kosher Food ➞ Vegan",
	},
	map[string]any{
		"value":    50,
		"label":    "ベジタリアン・ヴィーガン・ハラール・コーシャ ➞ ハラール",
		"label_en": "Vegitarian, Vegan, Halal and Kosher Food ➞ Halal",
	},
	map[string]any{
		"value":    51,
		"label":    "ベジタリアン・ヴィーガン・ハラール・コーシャ ➞ コーシャー",
		"label_en": "Vegitarian, Vegan, Halal and Kosher Food ➞ Kosher",
	},
	map[string]any{
		"value":    52,
		"label":    "肉・鶏肉 ➞ 乾燥肉",
		"label_en": "Meat & Poultry ➞ Dried Meat",
	},
	map[string]any{
		"value":    53,
		"label":    "肉・鶏肉 ➞ 内臓肉",
		"label_en": "Meat & Poultry ➞ Organ Meats & Offal",
	},
	map[string]any{
		"value":    54,
		"label":    "肉・鶏肉 ➞ 鶏肉・羽毛肉",
		"label_en": "Meat & Poultry ➞ Poultry & Feathered Game",
	},
	map[string]any{
		"value":    55,
		"label":    "肉・鶏肉 ➞ 加工肉",
		"label_en": "Meat & Poultry ➞ Processed Meat",
	},
	map[string]any{
		"value":    56,
		"label":    "肉・鶏肉 ➞ 未加工肉",
		"label_en": "Meat & Poultry ➞ Unprocessed Meat",
	},
	map[string]any{
		"value":    57,
		"label":    "オーガニック製品 ➞ オーガニック製品",
		"label_en": "Organic Products ➞ Organic Products",
	},
	map[string]any{"value": 58, "label": "水産物 ➞ 水産物", "label_en": "Seafood ➞ Seafood"},
	map[string]any{
		"value":    126,
		"label":    "水産物 ➞ 水産加工品",
		"label_en": "Seafood ➞ Processed seafood products",
	},
	map[string]any{
		"value":    59,
		"label":    "スナック ➞ 焼き菓子",
		"label_en": "Snacks ➞ Savoury Baked Products",
	},
	map[string]any{
		"value":    60,
		"label":    "スナック ➞ スナック菓子",
		"label_en": "Snacks ➞ Savoury snacks",
	},
	map[string]any{
		"value":    61,
		"label":    "スナック ➞ 野菜スナック",
		"label_en": "Snacks ➞ Vegetable snacks",
	},
	map[string]any{
		"value":    62,
		"label":    "スナック ➞ シリアル／フルーツバー",
		"label_en": "Snacks ➞ Cereals / Fruit Bars",
	},
	map[string]any{"value": 63, "label": "スナック ➞ クラッカー", "label_en": "Snacks ➞ Crackers"},
	map[string]any{
		"value":    64,
		"label":    "スナック ➞ デーツ＆ナツメヤシ製品",
		"label_en": "Snacks ➞ Dates & Date Palm Products",
	},
	map[string]any{
		"value":    65,
		"label":    "スナック ➞ ドライフルーツ＆野菜",
		"label_en": "Snacks ➞ Dried Fruit & Vegetables",
	},
	map[string]any{
		"value":    66,
		"label":    "スナック ➞ フィンガーフード",
		"label_en": "Snacks ➞ Finger Food",
	},
	map[string]any{"value": 67, "label": "スナック ➞ ナッツ類", "label_en": "Snacks ➞ Nuts"},
	map[string]any{"value": 68, "label": "スナック ➞ ポップコーン", "label_en": "Snacks ➞ Popcorn"},
	map[string]any{
		"value":    69,
		"label":    "スナック ➞ 塩味スナック",
		"label_en": "Snacks ➞ Salty Snacks",
	},
	map[string]any{
		"value":    107,
		"label":    "スナック ➞ せんべい、ヘルシー・スナック",
		"label_en": "Snacks ➞ Rice Crackers & Healthy Snacks",
	},
	map[string]any{
		"value":    71,
		"label":    "スプレッド、蜂蜜、ジャム ➞ スプレッド、蜂蜜、ジャム",
		"label_en": "Spreads, Honey & Jams ➞ Spreads, Honey & Jams",
	},
	map[string]any{
		"value":    72,
		"label":    "植物油脂・動物油脂 ➞ 植物油脂・動物油脂",
		"label_en": "Vegetable & Animal Oils & Fats ➞ Vegetable & Animal Oils & Fats",
	},
	map[string]any{
		"value":    73,
		"label":    "食品素材・原料 ➞ シリアル",
		"label_en": "Food Ingredients ➞ Breakfast Cereals",
	},
	map[string]any{
		"value":    74,
		"label":    "食品素材・原料 ➞ コーヒー豆",
		"label_en": "Food Ingredients ➞ Coffee Beans",
	},
	map[string]any{
		"value":    75,
		"label":    "食品素材・原料 ➞ 豆類",
		"label_en": "Food Ingredients ➞ Pulses",
	},
	map[string]any{
		"value":    76,
		"label":    "食品素材・原料 ➞ 米",
		"label_en": "Food Ingredients ➞ Rice",
	},
	map[string]any{
		"value":    77,
		"label":    "食品素材・原料 ➞ 小麦",
		"label_en": "Food Ingredients ➞ Wheat",
	},
	map[string]any{
		"value":    78,
		"label":    "包装資材・容器 ➞ 包装資材・容器",
		"label_en": "Packaging ➞ Packaging",
	},
	map[string]any{
		"value":    79,
		"label":    "衛生対策・品質管理 ➞ 衛生対策・品質管理",
		"label_en": "Sanitary ➞ Sanitary",
	},
	map[string]any{
		"value":    80,
		"label":    "食品機器関連 ➞ 調理・焼成機器",
		"label_en": "Food Equipment ➞ Cooking and Baking Equipment",
	},
	map[string]any{
		"value":    81,
		"label":    "食品機器関連 ➞ 冷蔵・冷凍機器",
		"label_en": "Food Equipment ➞ Refrigeration and Cooling Equipment",
	},
	map[string]any{
		"value":    82,
		"label":    "食品機器関連 ➞ 調理用品",
		"label_en": "Food Equipment ➞ Food Preparation Equipment",
	},
	map[string]any{
		"value":    83,
		"label":    "食品機器関連 ➞ 飲料機器",
		"label_en": "Food Equipment ➞ Beverage Equipment",
	},
	map[string]any{
		"value":    127,
		"label":    "食品機器関連 ➞ 食品製造機器",
		"label_en": "Food Equipment ➞ Food Manufacturing Equipment",
	},
	map[string]any{
		"value":    84,
		"label":    "情報・サービス ➞ 情報・サービス",
		"label_en": "Professional Services ➞ Professional Services",
	},
	map[string]any{
		"value":    128,
		"label":    "情報・サービス ➞ DXソリューション",
		"label_en": "Professional Services ➞ DX Solutuion",
	},
	map[string]any{
		"value":    85,
		"label":    "ワイン ➞ スティルワイン - 赤",
		"label_en": "Wine ➞ Still wines - red",
	},
	map[string]any{
		"value":    86,
		"label":    "ワイン ➞ スティルワイン - ロゼ",
		"label_en": "Wine ➞ Still wines - rose",
	},
	map[string]any{
		"value":    87,
		"label":    "ワイン ➞ スティルワイン - 白",
		"label_en": "Wine ➞ Still wines - white",
	},
	map[string]any{
		"value":    88,
		"label":    "ワイン ➞ スパークリングワイン",
		"label_en": "Wine ➞ Sparkling wines",
	},
	map[string]any{
		"value":    89,
		"label":    "ワイン ➞ フォーティファイドワイン",
		"label_en": "Wine ➞ Fortified wines",
	},
	map[string]any{
		"value":    90,
		"label":    "ワイン ➞ フレーバードワイン",
		"label_en": "Wine ➞ Flavored wines",
	},
	map[string]any{
		"value":    91,
		"label":    "ワイン ➞ デザートワイン",
		"label_en": "Wine ➞ Dessert wines",
	},
	map[string]any{
		"value":    92,
		"label":    "ワイン ➞ オーガニックワイン",
		"label_en": "Wine ➞ Organic wines",
	},
	map[string]any{
		"value":    93,
		"label":    "ワイン ➞ ヴィーガンワイン",
		"label_en": "Wine ➞ Vegan wines",
	},
	map[string]any{
		"value":    94,
		"label":    "ワイン ➞ フルーツワイン",
		"label_en": "Wine ➞ Fruit and flavored",
	},
	map[string]any{
		"value":    95,
		"label":    "ワイン ➞ ノンアルコールワイン",
		"label_en": "Wine ➞ Non-alcoholic wines",
	},
	map[string]any{
		"value":    96,
		"label":    "ワイン ➞ 低アルコールワイン",
		"label_en": "Wine ➞ Low-alcoholic wines",
	},
	map[string]any{
		"value":    97,
		"label":    "アルコール ➞ ビール",
		"label_en": "Alcoholic Beverage ➞ Beer",
	},
	map[string]any{
		"value":    98,
		"label":    "アルコール ➞ 日本酒",
		"label_en": "Alcoholic Beverage ➞ Japanese Sake",
	},
	map[string]any{
		"value":    99,
		"label":    "アルコール ➞ 焼酎",
		"label_en": "Alcoholic Beverage ➞ Shochu",
	},
	map[string]any{
		"value":    100,
		"label":    "アルコール ➞ フルーツリキュール",
		"label_en": "Alcoholic Beverage ➞ Fruit Liqueur",
	},
	map[string]any{
		"value":    101,
		"label":    "アルコール ➞ スピリッツ",
		"label_en": "Alcoholic Beverage ➞ Spirits",
	},
	map[string]any{
		"value":    102,
		"label":    "アルコール ➞ ウィスキー",
		"label_en": "Alcoholic Beverage ➞ Whiskey",
	},
	map[string]any{
		"value":    103,
		"label":    "アルコール ➞ その他",
		"label_en": "Alcoholic Beverage ➞ Others",
	},
	map[string]any{
		"value":    109,
		"label":    "AI 関連サービス ➞ 生産管理サービス",
		"label_en": "AI Services ➞ Production Management",
	},
	map[string]any{
		"value":    110,
		"label":    "AI 関連サービス ➞ 受発注サービス",
		"label_en": "AI Services ➞ Order Management",
	},
	map[string]any{
		"value":    111,
		"label":    "AI 関連サービス ➞ 画像検査",
		"label_en": "AI Services ➞ Image Inspection",
	},
	map[string]any{
		"value":    112,
		"label":    "AI 関連サービス ➞ 音声入力",
		"label_en": "AI Services ➞ Voice Input",
	},
	map[string]any{
		"value":    113,
		"label":    "AI 関連サービス ➞ 設備保全",
		"label_en": "AI Services ➞ Equipment Maintenance",
	},
	map[string]any{
		"value":    114,
		"label":    "AI 関連サービス ➞ 生成AI",
		"label_en": "AI Services ➞ Generative AI",
	},
	map[string]any{
		"value":    115,
		"label":    "AI 関連サービス ➞ データ分析",
		"label_en": "AI Services ➞ Data Analysis",
	},
	map[string]any{
		"value":    116,
		"label":    "AI 関連サービス ➞ 販促・プロモーション",
		"label_en": "AI Services ➞ Sales Promotion",
	},
	map[string]any{
		"value":    117,
		"label":    "AI 関連サービス ➞ その他",
		"label_en": "AI Services ➞ Others",
	},
	map[string]any{
		"value":    118,
		"label":    "物流サービス ➞ 総合物流サービス",
		"label_en": "Logistics ➞ Comprehensive Logistics",
	},
	map[string]any{
		"value":    119,
		"label":    "物流サービス ➞ 国内物流",
		"label_en": "Logistics ➞ Japan Domestic Logistics",
	},
	map[string]any{
		"value":    120,
		"label":    "物流サービス ➞ 輸入サポート",
		"label_en": "Logistics ➞ Import to Japan Support",
	},
	map[string]any{
		"value":    121,
		"label":    "物流サービス ➞ 輸出サポート",
		"label_en": "Logistics ➞ Export from Japan Support",
	},
	map[string]any{
		"value":    122,
		"label":    "物流サービス ➞ 倉庫",
		"label_en": "Logistics ➞ Warehousing",
	},
	map[string]any{"value": 123, "label": "物流サービス ➞ その他", "label_en": "Logistics ➞ Others"},
}

var TOPIC_VISITOR_INTERESTED_CATEGORIES_HCJ = []map[string]any{
	map[string]any{
		"value":    1,
		"label":    "厨房設備・調理機器 ➞ 自動ロボット(調理、配膳、清掃 等)",
		"label_en": "Food Service Equipment ➞ Automated Robots (Cooking, Serving, Cleaning, etc.)",
	},
	map[string]any{
		"value":    2,
		"label":    "厨房設備・調理機器 ➞ 加熱機器・熱調理機器",
		"label_en": "Food Service Equipment ➞ Heating & Thermal Cooking Equipment",
	},
	map[string]any{
		"value":    3,
		"label":    "厨房設備・調理機器 ➞ 冷凍・冷蔵機器(真空冷却機、フリーザー等)",
		"label_en": "Food Service Equipment ➞ Refrigeration & Freezing Equipment (vacuum coolers, freezers, etc.)",
	},
	map[string]any{
		"value":    4,
		"label":    "厨房設備・調理機器 ➞ 自動調理器(ブレンダー・真空包装機等)",
		"label_en": "Food Service Equipment ➞ Automatic Cooker (blenders, vacuum sealers, etc.)",
	},
	map[string]any{
		"value":    5,
		"label":    "厨房設備・調理機器 ➞ 調理道具(包丁・まな板・スライサー・鍋 等)",
		"label_en": "Food Service Equipment ➞ Cooking Utensils (Knives, Cutting Boards, Slicers, Pots, etc.)",
	},
	map[string]any{
		"value":    6,
		"label":    "厨房設備・調理機器 ➞ 厨房システム(温度管理、鮮度モニタリング等)",
		"label_en": "Food Service Equipment ➞ Kitchen Systems (temperature control, freshness monitoring, etc.)",
	},
	map[string]any{
		"value":    7,
		"label":    "厨房設備・調理機器 ➞ 換気扇",
		"label_en": "Food Service Equipment ➞ Ventilation Fans",
	},
	map[string]any{
		"value":    8,
		"label":    "厨房設備・調理機器 ➞ 洗浄機(食器、グラス、器具 等)",
		"label_en": "Food Service Equipment ➞ Washing Equipment (Dishwashers, Glass Washers, Utensil Washers, etc.)",
	},
	map[string]any{
		"value":    9,
		"label":    "厨房設備・調理機器 ➞ ごみ処理機",
		"label_en": "Food Service Equipment ➞ Waste Disposal",
	},
	map[string]any{
		"value":    10,
		"label":    "厨房設備・調理機器 ➞ ろ過器（水・油等）",
		"label_en": "Food Service Equipment ➞ Filters (water, oil, etc.)",
	},
	map[string]any{
		"value":    11,
		"label":    "厨房設備・調理機器 ➞ ピザオーブン、ピザ成型機",
		"label_en": "Food Service Equipment ➞ Pizza Ovens, Pizza Forming Machines",
	},
	map[string]any{
		"value":    12,
		"label":    "厨房設備・調理機器 ➞ ユニフォーム及びユニフォーム関連雑貨",
		"label_en": "Food Service Equipment ➞ Uniforms and Uniform Accessories",
	},
	map[string]any{
		"value":    13,
		"label":    "厨房設備・調理機器 ➞ システム(順番受付、モバイルオーダー等)",
		"label_en": "Food Service Equipment ➞ Systems (Queue Management, Mobile Ordering, etc.)",
	},
	map[string]any{
		"value":    14,
		"label":    "厨房設備・調理機器 ➞ タッチパネル券売機、飲食店向けPOSレジ",
		"label_en": "Food Service Equipment ➞ Touch Panel Ticket Vending Machines, POS Registers for Restaurants",
	},
	map[string]any{
		"value":    15,
		"label":    "厨房設備・調理機器 ➞ キッチンディスプレイ、キッチンプリンター",
		"label_en": "Food Service Equipment ➞ Kitchen Displays, Kitchen Printers",
	},
	map[string]any{
		"value":    16,
		"label":    "厨房設備・調理機器 ➞ その他",
		"label_en": "Food Service Equipment ➞ Food Service Equipment Others",
	},
	map[string]any{
		"value":    17,
		"label":    "ビュッフェ・サーブ・ショーケース ➞ ビュッフェ関連（保温器・トング・消耗品等）",
		"label_en": "Buffet / Serving / Showcase ➞ Buffet Equipment (warmers, tongs, consumables, etc.)",
	},
	map[string]any{
		"value":    18,
		"label":    "ビュッフェ・サーブ・ショーケース ➞ サーブ関連（ワゴン・大皿・トング 等）",
		"label_en": "Buffet / Serving / Showcase ➞ Serving Equipment (carts, large plates, tongs, etc.)",
	},
	map[string]any{
		"value":    19,
		"label":    "ビュッフェ・サーブ・ショーケース ➞ ショーケース（保温・保冷 等）",
		"label_en": "Buffet / Serving / Showcase ➞ Showcases (warming, cooling, etc.)",
	},
	map[string]any{
		"value":    20,
		"label":    "ビュッフェ・サーブ・ショーケース ➞ その他",
		"label_en": "Buffet / Serving / Showcase ➞ Buffet / Serving / Showcase Others",
	},
	map[string]any{
		"value":    21,
		"label":    "テーブルウェア・容器（使い捨て以外） ➞ 食器（陶磁器・漆器・銀器、ガラス器、竹製、木製等）",
		"label_en": "Tableware / Food Container (Non-disposable) ➞ Tableware (Ceramics, Lacquerware, Silverware, Glassware, Bamboo, Wooden, etc.)",
	},
	map[string]any{
		"value":    22,
		"label":    "テーブルウェア・容器（使い捨て以外） ➞ グラス",
		"label_en": "Tableware / Food Container (Non-disposable) ➞ Glassware",
	},
	map[string]any{
		"value":    23,
		"label":    "テーブルウェア・容器（使い捨て以外） ➞ バスケット、トレイ",
		"label_en": "Tableware / Food Container (Non-disposable) ➞ Baskets, Trays",
	},
	map[string]any{
		"value":    24,
		"label":    "テーブルウェア・容器（使い捨て以外） ➞ メラミン食器",
		"label_en": "Tableware / Food Container (Non-disposable) ➞ Melamine Tableware",
	},
	map[string]any{
		"value":    25,
		"label":    "テーブルウェア・容器（使い捨て以外） ➞ フードコンテナ・弁当箱（使い捨て以外）",
		"label_en": "Tableware / Food Container (Non-disposable) ➞ Food Containers & Bento Boxes (non-disposable)",
	},
	map[string]any{
		"value":    26,
		"label":    "テーブルウェア・容器（使い捨て以外） ➞ カトラリー・トング 等",
		"label_en": "Tableware / Food Container (Non-disposable) ➞ Cutlery & Tongs",
	},
	map[string]any{
		"value":    27,
		"label":    "テーブルウェア・容器（使い捨て以外） ➞ テーブルリネン・ランチョンマット・花器",
		"label_en": "Tableware / Food Container (Non-disposable) ➞ Table Linens, Placemats, Vases",
	},
	map[string]any{
		"value":    28,
		"label":    "テーブルウェア・容器（使い捨て以外） ➞ 卓上品（メニュー・調味料入等）",
		"label_en": "Tableware / Food Container (Non-disposable) ➞ Tabletop Items (menus, condiment holders, etc.)",
	},
	map[string]any{
		"value":    29,
		"label":    "テーブルウェア・容器（使い捨て以外） ➞ 卓上品（おしぼり・ナプキン等）",
		"label_en": "Tableware / Food Container (Non-disposable) ➞ Tabletop Items (wet towels, napkins, etc.)",
	},
	map[string]any{
		"value":    30,
		"label":    "テーブルウェア・容器（使い捨て以外） ➞ SDGs関連",
		"label_en": "Tableware / Food Container (Non-disposable) ➞ SDGs-related Products",
	},
	map[string]any{
		"value":    31,
		"label":    "テーブルウェア・容器（使い捨て以外） ➞ ワイン関連(グラス、アクセサリー、ソムリエナイフ、ワインセラー 等)",
		"label_en": "Tableware / Food Container (Non-disposable) ➞ Wine Accessories (Glasses, Accessories, Sommelier Knives, Wine Cellars, etc.)",
	},
	map[string]any{
		"value":    32,
		"label":    "テーブルウェア・容器（使い捨て以外） ➞ キャンドル(アロマ、LED 等)",
		"label_en": "Tableware / Food Container (Non-disposable) ➞ Candles (Aroma, LED, etc.)",
	},
	map[string]any{
		"value":    33,
		"label":    "テーブルウェア・容器（使い捨て以外） ➞ その他",
		"label_en": "Tableware / Food Container (Non-disposable) ➞ Tableware / Food Container (Non-disposable) Others",
	},
	map[string]any{
		"value":    34,
		"label":    "容器・包装（使い捨て） ➞ 食品容器・弁当箱（使い捨て）",
		"label_en": "Food Container / Packaging (Disposable) ➞ Food Containers & Bento Boxes (disposable)",
	},
	map[string]any{
		"value":    35,
		"label":    "容器・包装（使い捨て） ➞ カトラリー・箸・ストロー（使い捨て）",
		"label_en": "Food Container / Packaging (Disposable) ➞ Cutlery, Chopsticks, Straws (disposable)",
	},
	map[string]any{
		"value":    36,
		"label":    "容器・包装（使い捨て） ➞ パッケージ製品（袋 等）",
		"label_en": "Food Container / Packaging (Disposable) ➞ Packaging Products（bags, etc.）",
	},
	map[string]any{
		"value":    37,
		"label":    "容器・包装（使い捨て） ➞ パッケージ関連（デザイン加工サービス 等）",
		"label_en": "Food Container / Packaging (Disposable) ➞ Package Related（design processing, etc.）",
	},
	map[string]any{
		"value":    38,
		"label":    "容器・包装（使い捨て） ➞ おしぼり・ナプキン（使い捨て）",
		"label_en": "Food Container / Packaging (Disposable) ➞ Wet Towels & Napkins (disposable)",
	},
	map[string]any{
		"value":    39,
		"label":    "容器・包装（使い捨て） ➞ SDGs関連",
		"label_en": "Food Container / Packaging (Disposable) ➞ SDGs-related Products",
	},
	map[string]any{
		"value":    40,
		"label":    "容器・包装（使い捨て） ➞ その他",
		"label_en": "Food Container / Packaging (Disposable) ➞ Food Container / Packaging (Disposable) Others",
	},
	map[string]any{
		"value":    41,
		"label":    "食品・飲料・調味料 ➞ 生鮮食品（農産物・食肉・水産物等）",
		"label_en": "Food / Beverage / Seasoning ➞ Fresh Foods (produce, meat, seafood, etc.)",
	},
	map[string]any{
		"value":    42,
		"label":    "食品・飲料・調味料 ➞ 冷蔵食品",
		"label_en": "Food / Beverage / Seasoning ➞ Refrigerated Foods",
	},
	map[string]any{
		"value":    43,
		"label":    "食品・飲料・調味料 ➞ 冷凍食品",
		"label_en": "Food / Beverage / Seasoning ➞ Frozen Foods",
	},
	map[string]any{
		"value":    44,
		"label":    "食品・飲料・調味料 ➞ 加工食品（カット・半加工・完全調理等）",
		"label_en": "Food / Beverage / Seasoning ➞ Processed Foods (pre-cut, semi-prepared, fully cooked, etc.)",
	},
	map[string]any{
		"value":    45,
		"label":    "食品・飲料・調味料 ➞ 食品素材（小麦・米・豆・コーヒー豆 等）",
		"label_en": "Food / Beverage / Seasoning ➞ Food Ingredients (wheat, rice, beans, coffee beans, etc.)",
	},
	map[string]any{
		"value":    46,
		"label":    "食品・飲料・調味料 ➞ 飲料（酒類・水・茶 等）",
		"label_en": "Food / Beverage / Seasoning ➞ Beverages (Alcohol, Water, Tea, etc.)",
	},
	map[string]any{
		"value":    47,
		"label":    "食品・飲料・調味料 ➞ 菓子・べーカリー・デザート",
		"label_en": "Food / Beverage / Seasoning ➞ Sweets, Bakery, Desserts",
	},
	map[string]any{
		"value":    48,
		"label":    "食品・飲料・調味料 ➞ 乾物（鰹節・パスタ・缶詰 等）",
		"label_en": "Food / Beverage / Seasoning ➞ Dried Foods (dried bonito, pasta, canned goods, etc.)",
	},
	map[string]any{
		"value":    49,
		"label":    "食品・飲料・調味料 ➞ 調味料・油",
		"label_en": "Food / Beverage / Seasoning ➞ Seasonings & Oils",
	},
	map[string]any{
		"value":    50,
		"label":    "食品・飲料・調味料 ➞ オーガニック製品",
		"label_en": "Food / Beverage / Seasoning ➞ Organic Products",
	},
	map[string]any{
		"value":    51,
		"label":    "食品・飲料・調味料 ➞ ベジタリアン・ヴィーガン・ハラール・コーシャ食品",
		"label_en": "Food / Beverage / Seasoning ➞ Vegetarian, Vegan, Halal, Kosher",
	},
	map[string]any{
		"value":    52,
		"label":    "食品・飲料・調味料 ➞ JGAP・グローバルGAP認証品",
		"label_en": "Food / Beverage / Seasoning ➞ JGAP/Global GAP Certified Products",
	},
	map[string]any{
		"value":    53,
		"label":    "食品・飲料・調味料 ➞ 地域特産物",
		"label_en": "Food / Beverage / Seasoning ➞ Local Specialty Products",
	},
	map[string]any{
		"value":    54,
		"label":    "食品・飲料・調味料 ➞ その他",
		"label_en": "Food / Beverage / Seasoning ➞ Food / Beverage / Seasoning Others",
	},
	map[string]any{
		"value":    55,
		"label":    "おみやげ・地域産品 ➞ おみやげ",
		"label_en": "Souvenir / Local Product ➞ Souvenirs",
	},
	map[string]any{
		"value":    56,
		"label":    "おみやげ・地域産品 ➞ 地域産品",
		"label_en": "Souvenir / Local Product ➞ Local Products",
	},
	map[string]any{
		"value":    57,
		"label":    "おみやげ・地域産品 ➞ その他",
		"label_en": "Souvenir / Local Product ➞ Souvenir / Local Product Others",
	},
	map[string]any{
		"value":    58,
		"label":    "カフェ・ベーカリー・デザート ➞ 飲料（水、茶、コーヒー等）",
		"label_en": "Cafe / Bakery / Dessert ➞ Beverages (water, tea, coffee, etc.)",
	},
	map[string]any{
		"value":    59,
		"label":    "カフェ・ベーカリー・デザート ➞ 食材（素材・半加工・完全調理、冷凍 等）",
		"label_en": "Cafe / Bakery / Dessert ➞ Food Ingredients (material, semi-processed, fully cooked, etc.)",
	},
	map[string]any{
		"value":    60,
		"label":    "カフェ・ベーカリー・デザート ➞ 飲料関連器具（コーヒーマシン・カップウォーマー等）",
		"label_en": "Cafe / Bakery / Dessert ➞ Beverage-related Equipment (coffee machines, cup warmers, etc.)",
	},
	map[string]any{
		"value":    61,
		"label":    "カフェ・ベーカリー・デザート ➞ 調理設備・機器（菓子・製パン・軽食関連等）",
		"label_en": "Cafe / Bakery / Dessert ➞ Cooking Equipment and Appliances (for pastry, baking, light meals, etc.)",
	},
	map[string]any{
		"value":    62,
		"label":    "カフェ・ベーカリー・デザート ➞ ショーケース（保温・保冷 等）",
		"label_en": "Cafe / Bakery / Dessert ➞ Showcases (warming, cooling, etc.)",
	},
	map[string]any{
		"value":    63,
		"label":    "カフェ・ベーカリー・デザート ➞ テーブルウェア（食器・カトラリー メニュー等）",
		"label_en": "Cafe / Bakery / Dessert ➞ Tableware (dishes, cutlery, menus, etc.)",
	},
	map[string]any{
		"value":    64,
		"label":    "カフェ・ベーカリー・デザート ➞ 容器・弁当箱（使い捨て）",
		"label_en": "Cafe / Bakery / Dessert ➞ Containers & Bento Boxes (disposable)",
	},
	map[string]any{
		"value":    65,
		"label":    "カフェ・ベーカリー・デザート ➞ パッケージ関連",
		"label_en": "Cafe / Bakery / Dessert ➞ Packaging-related",
	},
	map[string]any{
		"value":    66,
		"label":    "カフェ・ベーカリー・デザート ➞ システム関連（サイネージ・ＰＯＳレジ等）",
		"label_en": "Cafe / Bakery / Dessert ➞ System-related (digital signage, POS registers, etc.)",
	},
	map[string]any{
		"value":    67,
		"label":    "カフェ・ベーカリー・デザート ➞ アイスクリーム・ソルベ",
		"label_en": "Cafe / Bakery / Dessert ➞ Ice Cream and Sorbet",
	},
	map[string]any{
		"value":    68,
		"label":    "カフェ・ベーカリー・デザート ➞ 書籍・雑誌",
		"label_en": "Cafe / Bakery / Dessert ➞ Cafe / Bakery / Dessert",
	},
	map[string]any{
		"value":    69,
		"label":    "カフェ・ベーカリー・デザート ➞ その他",
		"label_en": "Cafe / Bakery / Dessert ➞ Cafe / Bakery / Dessert Others",
	},
	map[string]any{
		"value":    70,
		"label":    "衛生・清掃・HACCP ➞ 洗浄消毒装置・殺菌装置",
		"label_en": "Sanitation / Cleaning / HACCP ➞ Washing and Disinfection Equipment, Sterilization Equipment",
	},
	map[string]any{
		"value":    71,
		"label":    "衛生・清掃・HACCP ➞ 消毒剤・洗浄剤",
		"label_en": "Sanitation / Cleaning / HACCP ➞ Disinfectants, Detergents",
	},
	map[string]any{
		"value":    72,
		"label":    "衛生・清掃・HACCP ➞ 清掃用品",
		"label_en": "Sanitation / Cleaning / HACCP ➞ Cleaning Supplies",
	},
	map[string]any{
		"value":    73,
		"label":    "衛生・清掃・HACCP ➞ 洗濯・乾燥装置",
		"label_en": "Sanitation / Cleaning / HACCP ➞ Drying Equipment",
	},
	map[string]any{
		"value":    74,
		"label":    "衛生・清掃・HACCP ➞ 衛生衣料・手袋",
		"label_en": "Sanitation / Cleaning / HACCP ➞ Sanitary Clothing & Gloves",
	},
	map[string]any{
		"value":    75,
		"label":    "衛生・清掃・HACCP ➞ 消臭・フレグランス・分煙",
		"label_en": "Sanitation / Cleaning / HACCP ➞ Deodorants, Fragrances, Smoking Separation",
	},
	map[string]any{
		"value":    76,
		"label":    "衛生・清掃・HACCP ➞ 防虫防鼠・害虫害獣駆除",
		"label_en": "Sanitation / Cleaning / HACCP ➞ Pest Control (insects, rodents, etc.)",
	},
	map[string]any{
		"value":    77,
		"label":    "衛生・清掃・HACCP ➞ 人材サービス・人材育成・清掃マニュアル",
		"label_en": "Sanitation / Cleaning / HACCP ➞ Labor Saving Services, Human Resource Development",
	},
	map[string]any{
		"value":    78,
		"label":    "衛生・清掃・HACCP ➞ 清掃システム・清掃カメラ・清掃ロボット",
		"label_en": "Sanitation / Cleaning / HACCP ➞ Books and Magazines",
	},
	map[string]any{
		"value":    79,
		"label":    "衛生・清掃・HACCP ➞ その他",
		"label_en": "Sanitation / Cleaning / HACCP ➞ Sanitation / Cleaning / HACCP Others",
	},
	map[string]any{
		"value":    80,
		"label":    "設計改修・開業支援 ➞ リノベーション・リフォーム",
		"label_en": "Design Refurbishment / Opening Support ➞ Renovation, Remodeling",
	},
	map[string]any{
		"value":    81,
		"label":    "設計改修・開業支援 ➞ 空間・インテリア・家具デザイン",
		"label_en": "Design Refurbishment / Opening Support ➞ Space, Interior, Furniture Design",
	},
	map[string]any{
		"value":    82,
		"label":    "設計改修・開業支援 ➞ 資材（壁材・床材 等）",
		"label_en": "Design Refurbishment / Opening Support ➞ Materials (walls, floors, etc.)",
	},
	map[string]any{
		"value":    83,
		"label":    "設計改修・開業支援 ➞ 開業支援・コンサルティング・フランチャイズ",
		"label_en": "Design Refurbishment / Opening Support ➞ Opening Support, Consulting, Franchising",
	},
	map[string]any{
		"value":    84,
		"label":    "設計改修・開業支援 ➞ その他",
		"label_en": "Design Refurbishment / Opening Support ➞ Design Refurbishment / Opening Support Others",
	},
	map[string]any{
		"value":    85,
		"label":    "温浴施設・サウナ ➞ 温浴設備・資材",
		"label_en": "Spa / Bath / Sauna ➞ Bathing Facilities & Supplies",
	},
	map[string]any{
		"value":    86,
		"label":    "温浴施設・サウナ ➞ サウナ設備・資材",
		"label_en": "Spa / Bath / Sauna ➞ Sauna Equipment & Supplies",
	},
	map[string]any{
		"value":    87,
		"label":    "温浴施設・サウナ ➞ 備品（タオル・マットレス 等）",
		"label_en": "Spa / Bath / Sauna ➞ Supplies (towels, mattresses, etc.)",
	},
	map[string]any{
		"value":    88,
		"label":    "温浴施設・サウナ ➞ 消耗品（シャンプー 等）",
		"label_en": "Spa / Bath / Sauna ➞ Consumables (shampoo, etc.)",
	},
	map[string]any{
		"value":    89,
		"label":    "温浴施設・サウナ ➞ スパ・エステティック機器",
		"label_en": "Spa / Bath / Sauna ➞ Spa & Aesthetic Equipment",
	},
	map[string]any{
		"value":    90,
		"label":    "温浴施設・サウナ ➞ 設計・資材",
		"label_en": "Spa / Bath / Sauna ➞ Design & Materials",
	},
	map[string]any{
		"value":    91,
		"label":    "温浴施設・サウナ ➞ テクノロジー（POS・決済・混雑状況 等）",
		"label_en": "Spa / Bath / Sauna ➞ Technology (POS, payment, congestion status, etc.)",
	},
	map[string]any{
		"value":    92,
		"label":    "温浴施設・サウナ ➞ その他",
		"label_en": "Spa / Bath / Sauna ➞ Spa / Bath / Sauna Others",
	},
	map[string]any{
		"value":    93,
		"label":    "エクステリア関連 ➞ 屋外家具・日よけ",
		"label_en": "Exterior ➞ Outdoor Furniture, Sunshades",
	},
	map[string]any{
		"value":    94,
		"label":    "エクステリア関連 ➞ 屋外用暖房・冷風機",
		"label_en": "Exterior ➞ Outdoor Heaters & Coolers",
	},
	map[string]any{
		"value":    95,
		"label":    "エクステリア関連 ➞ 屋外照明",
		"label_en": "Exterior ➞ Outdoor Lighting",
	},
	map[string]any{
		"value":    96,
		"label":    "エクステリア関連 ➞ アウトドア調理器具",
		"label_en": "Exterior ➞ Outdoor Cooking Equipment",
	},
	map[string]any{
		"value":    97,
		"label":    "エクステリア関連 ➞ その他",
		"label_en": "Exterior ➞ Exterior Others",
	},
	map[string]any{
		"value":    98,
		"label":    "インテリア（ホテル共用部・外食・給食） ➞ 家具(テーブル・椅子・ソファ 等)",
		"label_en": "Interior (Hotel Common Areas / Food Service / Meal Service) ➞ Furniture (tables, chairs, sofas, etc.)",
	},
	map[string]any{
		"value":    99,
		"label":    "インテリア（ホテル共用部・外食・給食） ➞ 照明器具",
		"label_en": "Interior (Hotel Common Areas / Food Service / Meal Service) ➞ Lighting Fixtures",
	},
	map[string]any{
		"value":    100,
		"label":    "インテリア（ホテル共用部・外食・給食） ➞ 水まわり（手洗 等）",
		"label_en": "Interior (Hotel Common Areas / Food Service / Meal Service) ➞ Wet areas (washbasins, etc.)",
	},
	map[string]any{
		"value":    101,
		"label":    "インテリア（ホテル共用部・外食・給食） ➞ カーペット・マット・ブラインド 等",
		"label_en": "Interior (Hotel Common Areas / Food Service / Meal Service) ➞ Carpets, Mats, Blinds, etc.",
	},
	map[string]any{
		"value":    102,
		"label":    "インテリア（ホテル共用部・外食・給食） ➞ AV機器（テレビ・モニター・スピーカー 等）",
		"label_en": "Interior (Hotel Common Areas / Food Service / Meal Service) ➞ AV Equipment (TVs, monitors, speakers, etc.)",
	},
	map[string]any{
		"value":    103,
		"label":    "インテリア（ホテル共用部・外食・給食） ➞ 装飾（アート・花器等）",
		"label_en": "Interior (Hotel Common Areas / Food Service / Meal Service) ➞ Decorations (art, vases, etc.)",
	},
	map[string]any{
		"value":    104,
		"label":    "インテリア（ホテル共用部・外食・給食） ➞ 設計資材（扉・壁材 等）",
		"label_en": "Interior (Hotel Common Areas / Food Service / Meal Service) ➞ Design Materials (doors, walls, etc.)",
	},
	map[string]any{
		"value":    105,
		"label":    "インテリア（ホテル共用部・外食・給食） ➞ その他",
		"label_en": "Interior (Hotel Common Areas / Food Service / Meal Service) ➞ Interior (Hotel Common Areas / Food Service / Meal Service) Others",
	},
	map[string]any{
		"value":    106,
		"label":    "客室備品・アメニティ ➞ 寝具（ベッド・枕 等）",
		"label_en": "Guest Room Equipment / Amenity ➞ Bedding (beds, pillows, etc.)",
	},
	map[string]any{
		"value":    107,
		"label":    "客室備品・アメニティ ➞ 家具(テーブル・椅子・ソファ 等)",
		"label_en": "Guest Room Equipment / Amenity ➞ Furniture (tables, chairs, sofas, etc.)",
	},
	map[string]any{
		"value":    108,
		"label":    "客室備品・アメニティ ➞ カーペット・マット・カーテン・ブラインド 等",
		"label_en": "Guest Room Equipment / Amenity ➞ Carpets, Mats, Curtains, Blinds, etc.",
	},
	map[string]any{
		"value":    109,
		"label":    "客室備品・アメニティ ➞ 照明器具",
		"label_en": "Guest Room Equipment / Amenity ➞ Lighting Fixtures",
	},
	map[string]any{
		"value":    110,
		"label":    "客室備品・アメニティ ➞ ファブリック（リネン・タオル・パジャマ 等)",
		"label_en": "Guest Room Equipment / Amenity ➞ Fabrics (linen, towels, pajamas, etc.)",
	},
	map[string]any{
		"value":    111,
		"label":    "客室備品・アメニティ ➞ アメニティ（歯ブラシ・シャンプー・化粧品等）",
		"label_en": "Guest Room Equipment / Amenity ➞ Amenities (toothbrushes, shampoo, cosmetics, etc.)",
	},
	map[string]any{
		"value":    112,
		"label":    "客室備品・アメニティ ➞ 客室備品（ハンガー・ゴミ箱・金庫 等）",
		"label_en": "Guest Room Equipment / Amenity ➞ Guest Room Supplies (hangers, trash cans, safes, etc.)",
	},
	map[string]any{
		"value":    113,
		"label":    "客室備品・アメニティ ➞ スリッパ・サンダル 等",
		"label_en": "Guest Room Equipment / Amenity ➞ Slippers, Sandals, etc.",
	},
	map[string]any{
		"value":    114,
		"label":    "客室備品・アメニティ ➞ AV機器（テレビ・電話 等）",
		"label_en": "Guest Room Equipment / Amenity ➞ AV Equipment (TVs, telephones, etc.)",
	},
	map[string]any{
		"value":    115,
		"label":    "客室備品・アメニティ ➞ 電気製品（ドライヤー・冷蔵庫・空気清浄機 等）",
		"label_en": "Guest Room Equipment / Amenity ➞ Electrical Appliances (hair dryers, refrigerators, air purifiers, etc.)",
	},
	map[string]any{
		"value":    116,
		"label":    "客室備品・アメニティ ➞ 消臭・フレグランス",
		"label_en": "Guest Room Equipment / Amenity ➞ Deodorants, Fragrances",
	},
	map[string]any{
		"value":    117,
		"label":    "客室備品・アメニティ ➞ 客室内用飲料（水・珈琲・お茶 等）",
		"label_en": "Guest Room Equipment / Amenity ➞ In-room Beverages (water, coffee, tea, etc.)",
	},
	map[string]any{
		"value":    118,
		"label":    "客室備品・アメニティ ➞ 貸出備品（ズボンプレッサー・充電器 等）",
		"label_en": "Guest Room Equipment / Amenity ➞ Rental Supplies (trouser press, chargers, etc.)",
	},
	map[string]any{
		"value":    119,
		"label":    "客室備品・アメニティ ➞ 設計資材（扉・壁材 等）",
		"label_en": "Guest Room Equipment / Amenity ➞ Design Materials (doors, wall materials, etc.)",
	},
	map[string]any{
		"value":    120,
		"label":    "客室備品・アメニティ ➞ SDGs関連商品",
		"label_en": "Guest Room Equipment / Amenity ➞ SDGs-related Products",
	},
	map[string]any{
		"value":    121,
		"label":    "客室備品・アメニティ ➞ その他",
		"label_en": "Guest Room Equipment / Amenity ➞ Guest Room Equipment / Amenity Others",
	},
	map[string]any{
		"value":    122,
		"label":    "テクノロジー（システム） ➞ ホテル管理システム（PMS）",
		"label_en": "Technology (System) ➞ Property Management System (PMS)",
	},
	map[string]any{
		"value":    123,
		"label":    "テクノロジー（システム） ➞ レベニューマネジメントシステム・ダイナミックプライシングシステム",
		"label_en": "Technology (System) ➞ Revenue Management System, Dynamic Pricing System",
	},
	map[string]any{
		"value":    124,
		"label":    "テクノロジー（システム） ➞ サイトコントローラー",
		"label_en": "Technology (System) ➞ Site Controller",
	},
	map[string]any{
		"value":    125,
		"label":    "テクノロジー（システム） ➞ 口コミ関連システム",
		"label_en": "Technology (System) ➞ Review Management System",
	},
	map[string]any{
		"value":    126,
		"label":    "テクノロジー（システム） ➞ 混雑状況関連システム",
		"label_en": "Technology (System) ➞ Congestion Status Management System",
	},
	map[string]any{
		"value":    127,
		"label":    "テクノロジー（システム） ➞ セルフオーダーシステム",
		"label_en": "Technology (System) ➞ Self-Ordering System",
	},
	map[string]any{
		"value":    128,
		"label":    "テクノロジー（システム） ➞ POSシステム",
		"label_en": "Technology (System) ➞ POS System",
	},
	map[string]any{
		"value":    129,
		"label":    "テクノロジー（システム） ➞ 人事関連システム・シフト管理システム、人材派遣システム",
		"label_en": "Technology (System) ➞ Human Resources Systems, Shift Management Systems, Staffing Systems",
	},
	map[string]any{
		"value":    130,
		"label":    "テクノロジー（システム） ➞ 無人関連（駐車場・チェックイン・コンビニ等）",
		"label_en": "Technology (System) ➞ Unmanned Solutions (parking, check-in, convenience stores, etc.)",
	},
	map[string]any{
		"value":    131,
		"label":    "テクノロジー（システム） ➞ デジタルチップシステム",
		"label_en": "Technology (System) ➞ Digital Tip System",
	},
	map[string]any{
		"value":    132,
		"label":    "テクノロジー（システム） ➞ 在庫管理システム",
		"label_en": "Technology (System) ➞ Inventory Management System",
	},
	map[string]any{
		"value":    133,
		"label":    "テクノロジー（システム） ➞ その他",
		"label_en": "Technology (System) ➞ Technology (System) Others",
	},
	map[string]any{
		"value":    134,
		"label":    "テクノロジー（端末・周辺機器） ➞ チェックイン端末",
		"label_en": "Technology (Terminal / Peripheral Equipment) ➞ Check-In Machine",
	},
	map[string]any{
		"value":    135,
		"label":    "テクノロジー（端末・周辺機器） ➞ POSレジ・決済端末",
		"label_en": "Technology (Terminal / Peripheral Equipment) ➞ POS Register, Payment Terminal",
	},
	map[string]any{
		"value":    136,
		"label":    "テクノロジー（端末・周辺機器） ➞ テレビ・モニター・電子棚札",
		"label_en": "Technology (Terminal / Peripheral Equipment) ➞ TV, Monitors, Electronic Shelf Labels",
	},
	map[string]any{
		"value":    137,
		"label":    "テクノロジー（端末・周辺機器） ➞ タブレット・スマートフォン・電話機",
		"label_en": "Technology (Terminal / Peripheral Equipment) ➞ Tablets, Smartphones, Telephones",
	},
	map[string]any{
		"value":    138,
		"label":    "テクノロジー（端末・周辺機器） ➞ スマートロック・周辺機器",
		"label_en": "Technology (Terminal / Peripheral Equipment) ➞ Smart Locks, Peripheral Devices",
	},
	map[string]any{
		"value":    139,
		"label":    "テクノロジー（端末・周辺機器） ➞ 防犯カメラ・監視カメラ",
		"label_en": "Technology (Terminal / Peripheral Equipment) ➞ Security Cameras, Surveillance Cameras",
	},
	map[string]any{
		"value":    140,
		"label":    "テクノロジー（端末・周辺機器） ➞ インカム・トランシーバー",
		"label_en": "Technology (Terminal / Peripheral Equipment) ➞ Intercoms, Walkie-Talkies",
	},
	map[string]any{
		"value":    141,
		"label":    "テクノロジー（端末・周辺機器） ➞ 通信関連",
		"label_en": "Technology (Terminal / Peripheral Equipment) ➞ Communication Systems",
	},
	map[string]any{
		"value":    142,
		"label":    "テクノロジー（端末・周辺機器） ➞ その他",
		"label_en": "Technology (Terminal / Peripheral Equipment) ➞ Technology (Terminal / Peripheral Equipment) Others",
	},
	map[string]any{
		"value":    143,
		"label":    "テクノロジー（動画・音楽コンテンツ） ➞ ホテルインフォメーション",
		"label_en": "Technology (Video / Music Contents) ➞ Hotel Information",
	},
	map[string]any{
		"value":    144,
		"label":    "テクノロジー（動画・音楽コンテンツ） ➞ 動画・音楽コンテンツ",
		"label_en": "Technology (Video / Music Contents) ➞ Video & Music Contents",
	},
	map[string]any{
		"value":    145,
		"label":    "テクノロジー（動画・音楽コンテンツ） ➞ その他",
		"label_en": "Technology (Video / Music Contents) ➞ Technology (Video / Music Contents) Others",
	},
	map[string]any{
		"value":    146,
		"label":    "テクノロジー（ロボット・AI・RPA） ➞ ロボット（配膳・配送・清掃・案内・警備 等）",
		"label_en": "Technology (Robot / AI / RPA) ➞ Robots (for serving, delivery, cleaning, guidance, security, etc.)",
	},
	map[string]any{
		"value":    147,
		"label":    "テクノロジー（ロボット・AI・RPA） ➞ 遠隔対応・アバター 等",
		"label_en": "Technology (Robot / AI / RPA) ➞ Remote Support, Avatars, etc.",
	},
	map[string]any{
		"value":    148,
		"label":    "テクノロジー（ロボット・AI・RPA） ➞ RPA・AI関連（チャットボット・通訳・多言語対応・AIエージェント 等）",
		"label_en": "Technology (Robot / AI / RPA) ➞ RPA and AI Solutions (Chatbots, Translation, Multilingual Support, AI Agents, etc.)",
	},
	map[string]any{
		"value":    149,
		"label":    "テクノロジー（ロボット・AI・RPA） ➞ その他",
		"label_en": "Technology (Robot / AI / RPA) ➞ Technology (Robot / AI / RPA) Others",
	},
	map[string]any{
		"value":    150,
		"label":    "テクノロジー（給食関連） ➞ 給食関連システム（栄養管理・給食管理・献立支援・精算等）",
		"label_en": "Technology (Meal Service related) ➞ Meal Management Systems (nutrition, meal service , menu support, billing, etc.)",
	},
	map[string]any{
		"value":    151,
		"label":    "テクノロジー（給食関連） ➞ その他",
		"label_en": "Technology (Meal Service related) ➞ Technology (Meal Service related) Others",
	},
	map[string]any{
		"value":    152,
		"label":    "テクノロジー（PR・マーケティング） ➞ 口コミ関連システム",
		"label_en": "Public Relation / Marketing ➞ Review Management System",
	},
	map[string]any{
		"value":    153,
		"label":    "テクノロジー（PR・マーケティング） ➞ データ分析（BIツール・クチコミ分析 等）",
		"label_en": "Public Relation / Marketing ➞ Data Analytics (BI tools, review analysis, etc.)",
	},
	map[string]any{
		"value":    154,
		"label":    "テクノロジー（PR・マーケティング） ➞ CRM・SFA・MA",
		"label_en": "Public Relation / Marketing ➞ CRM, SFA, MA",
	},
	map[string]any{
		"value":    155,
		"label":    "テクノロジー（PR・マーケティング） ➞ 販促サービス（広告・SNS運用代行など）・販促グッズ",
		"label_en": "Public Relation / Marketing ➞ Promotion Services (advertising, SNS management, etc.), Promotional Goods",
	},
	map[string]any{
		"value":    156,
		"label":    "テクノロジー（PR・マーケティング） ➞ チャットボット",
		"label_en": "Public Relation / Marketing ➞ Chatbots",
	},
	map[string]any{
		"value":    157,
		"label":    "テクノロジー（PR・マーケティング） ➞ WEB制作関連",
		"label_en": "Public Relation / Marketing ➞ Web Production related",
	},
	map[string]any{
		"value":    158,
		"label":    "テクノロジー（PR・マーケティング） ➞ コンサルティング",
		"label_en": "Public Relation / Marketing ➞ Consulting",
	},
	map[string]any{
		"value":    159,
		"label":    "テクノロジー（PR・マーケティング） ➞ その他",
		"label_en": "Public Relation / Marketing ➞ Public Relation / Marketing Others",
	},
	map[string]any{
		"value":    160,
		"label":    "人事ソリューション / 人材採用 / 従業員満足向上 ➞ 人材採用・外国人採用",
		"label_en": "HR Solution / Recruiting / Improving Employee Satisfaction ➞ Recruitment, Foreign Worker Recruitment",
	},
	map[string]any{
		"value":    161,
		"label":    "人事ソリューション / 人材採用 / 従業員満足向上 ➞ 人材派遣・アルバイト採用",
		"label_en": "HR Solution / Recruiting / Improving Employee Satisfaction ➞ Temporary Staffing, Part-Time Recruitment",
	},
	map[string]any{
		"value":    162,
		"label":    "人事ソリューション / 人材採用 / 従業員満足向上 ➞ 人材育成支援",
		"label_en": "HR Solution / Recruiting / Improving Employee Satisfaction ➞ Employee Training Support",
	},
	map[string]any{
		"value":    163,
		"label":    "人事ソリューション / 人材採用 / 従業員満足向上 ➞ 従業員満足度向上・働き方改革",
		"label_en": "HR Solution / Recruiting / Improving Employee Satisfaction ➞ Employee Satisfaction Improvement, Work Style Reform",
	},
	map[string]any{
		"value":    164,
		"label":    "人事ソリューション / 人材採用 / 従業員満足向上 ➞ 人事労務管理システム",
		"label_en": "HR Solution / Recruiting / Improving Employee Satisfaction ➞ HR and Labor Management System",
	},
	map[string]any{
		"value":    165,
		"label":    "人事ソリューション / 人材採用 / 従業員満足向上 ➞ その他",
		"label_en": "HR Solution / Recruiting / Improving Employee Satisfaction ➞ HR Solution / Recruiting / Improving Employee Satisfaction Others",
	},
	map[string]any{
		"value":    166,
		"label":    "省エネ・省コスト対策 ➞ 省エネ・省コスト製品",
		"label_en": "Energy / Cost-saving ➞ Energy Saving, Cost Reduction Products",
	},
	map[string]any{
		"value":    167,
		"label":    "省エネ・省コスト対策 ➞ 省エネ設計・コンサルティング関連",
		"label_en": "Energy / Cost-saving ➞ Energy Saving Design, Consulting",
	},
	map[string]any{
		"value":    168,
		"label":    "省エネ・省コスト対策 ➞ その他",
		"label_en": "Energy / Cost-saving ➞ Energy / Cost-saving Others",
	},
	map[string]any{
		"value":    169,
		"label":    "ペットツーリズム ➞ ペットフード・おやつ",
		"label_en": "Pet Tourism ➞ Pet Food and Treats",
	},
	map[string]any{
		"value":    170,
		"label":    "ペットツーリズム ➞ ペットグッズ・トイレタリー",
		"label_en": "Pet Tourism ➞ Pet Accessories and Toiletries",
	},
	map[string]any{
		"value":    171,
		"label":    "ペットツーリズム ➞ トリミング用品",
		"label_en": "Pet Tourism ➞ Grooming Supplies",
	},
	map[string]any{
		"value":    172,
		"label":    "ペットツーリズム ➞ ペット専用アメニティ",
		"label_en": "Pet Tourism ➞ Pet-specific Amenities",
	},
	map[string]any{
		"value":    173,
		"label":    "ペットツーリズム ➞ ドッグラン ・ ペット用プール製品",
		"label_en": "Pet Tourism ➞ Dog Run and Pet Pool Products",
	},
	map[string]any{
		"value":    174,
		"label":    "ペットツーリズム ➞ 環境機器(空気清浄機 ・ 脱臭機 等)",
		"label_en": "Pet Tourism ➞ Environmental Equipment (Air Purifiers, Deodorizers, etc.)",
	},
	map[string]any{
		"value":    175,
		"label":    "ペットツーリズム ➞ 消毒・除菌用品",
		"label_en": "Pet Tourism ➞ Disinfection and Sanitization Products",
	},
	map[string]any{
		"value":    176,
		"label":    "ペットツーリズム ➞ その他",
		"label_en": "Pet Tourism ➞ Pet Tourism and Other Services",
	},
	map[string]any{"value": 177, "label": "その他 ➞ その他", "label_en": "Others ➞ Others"},
}

var TOPIC_VISITOR_ZONE_OPTIONS_FOODEX = []map[string]any{
	map[string]any{"value": 1, "label": "FOODEX TECH", "label_en": "FOODEX TECH"},
	map[string]any{"value": 2, "label": "FOODEX WINE", "label_en": "FOODEX WINE"},
	map[string]any{
		"value":    3,
		"label":    "FOODEX WINE 日本ワインパビリオン",
		"label_en": "FOODEX WINE Japanese Wine Pavilion",
	},
	map[string]any{
		"value":    4,
		"label":    "FOODEX Frozen FRESH FOOD",
		"label_en": "FOODEX Frozen Fresh Products",
	},
	map[string]any{
		"value":    5,
		"label":    "FOODEX Frozen スイーツ",
		"label_en": "FOODEX Frozen Sweets & Snacks",
	},
	map[string]any{
		"value":    6,
		"label":    "FOODEX Frozen 加工食品",
		"label_en": "FOODEX Frozen Processed Food",
	},
	map[string]any{
		"value":    7,
		"label":    "FOODEX Frozen 機械技術",
		"label_en": "FOODEX Frozen Machines & Technology",
	},
	map[string]any{
		"value":    8,
		"label":    "FOODEX Frozen 日本の地域産品・世界のローカルフード",
		"label_en": "FOODEX Frozen Local Food from Worldwide",
	},
	map[string]any{"value": 9, "label": "オーガニック", "label_en": "Organic"},
	map[string]any{"value": 10, "label": "スイーツ＆スナック", "label_en": "Sweets & Snacks"},
	map[string]any{"value": 11, "label": "ドリンク＆アルコール", "label_en": "Beverage"},
	map[string]any{"value": 12, "label": "にっぽん食輸出展", "label_en": "Exports from Japan"},
	map[string]any{"value": 13, "label": "ヘルスケア", "label_en": "Healthcare Foods"},
	map[string]any{"value": 14, "label": "加工食品", "label_en": "Processed Food"},
	map[string]any{"value": 15, "label": "情報・サービス", "label_en": "Publication / Consultant"},
	map[string]any{"value": 16, "label": "食品安全対策展", "label_en": "Food Safety"},
	map[string]any{"value": 17, "label": "水産", "label_en": "Seafood"},
	map[string]any{"value": 18, "label": "全国食品博", "label_en": "Local Specialities in Japan"},
	map[string]any{
		"value":    19,
		"label":    "代替食・食品素材",
		"label_en": "Plant Based / Alternative Food",
	},
	map[string]any{"value": 20, "label": "畜産", "label_en": "Meat"},
	map[string]any{"value": 21, "label": "調味料", "label_en": "Condiments / Seasoning"},
	map[string]any{"value": 22, "label": "農産", "label_en": "Agricultural Food"},
	map[string]any{"value": 23, "label": "輸入食品", "label_en": "Imports Food"},
	map[string]any{"value": 24, "label": "海外ナショナルパビリオン", "label_en": "National Pavilion"},
	map[string]any{"value": 25, "label": "SAKE JAPAN", "label_en": "SAKE JAPAN"},
	map[string]any{"value": 26, "label": "物流", "label_en": "Logistics"},
	map[string]any{"value": 27, "label": "スタートアップ", "label_en": "Start up"},
	map[string]any{"value": 28, "label": "食×AI", "label_en": "Food × AI"},
	map[string]any{
		"value":    29,
		"label":    "ハラル・ヴィーガン・コーシャ",
		"label_en": "Halal,Vegan, Kosher",
	},
	map[string]any{"value": 30, "label": "ワールドフード", "label_en": "World Food"},
	map[string]any{"value": 31, "label": "輸出支援サービス", "label_en": "Export Support Service"},
	map[string]any{"value": 32, "label": "チルド・日配", "label_en": "Chilled Food"},
	map[string]any{"value": 99, "label": "食肉産業展", "label_en": "Japan Meat Industry Fair"},
}

var TOPIC_VISITOR_ZONE_OPTIONS_HCJ = []map[string]any{
	map[string]any{"value": 1, "label": "厨房設備・機器", "label_en": "Food Service Equipment"},
	map[string]any{
		"value":    2,
		"label":    "業務用食材・飲料・地域産品",
		"label_en": "Food & Drinks for Professional-Use / Local Products in Japan",
	},
	map[string]any{
		"value":    3,
		"label":    "カフェ・ベーカリー・デザート",
		"label_en": "Cafe / Bakery / Dessert",
	},
	map[string]any{"value": 4, "label": "テーブルウェア", "label_en": "Tableware"},
	map[string]any{
		"value":    5,
		"label":    "給食・弁当関連・包装資材容器",
		"label_en": "Food Container / Package",
	},
	map[string]any{
		"value":    6,
		"label":    "ホスピタリティデザイン東京/ 客室備品・アメニティ・家具・インテリア",
		"label_en": "Interior / Guest Room Equipment / Amenity",
	},
	map[string]any{
		"value":    7,
		"label":    "ホスピタリティデザイン東京/ ホテル・旅館・飲食店向け設計・改修",
		"label_en": "Design and Renovation for Hotels, Ryokans, Restaurants",
	},
	map[string]any{
		"value":    8,
		"label":    "ホスピタリティデザイン東京/ 屋外・エクステリア・レジャー",
		"label_en": "Outdoor Equipment / Exterior",
	},
	map[string]any{
		"value":    9,
		"label":    "エコ・省エネ対策",
		"label_en": "Ecology / Energy Saving Measures",
	},
	map[string]any{"value": 10, "label": "衛生・清掃", "label_en": "Sanitation / Cleaning"},
	map[string]any{
		"value":    11,
		"label":    "JAPAN サウナ・スパEXPO",
		"label_en": "Japan Sauna & Spa EXPO",
	},
	map[string]any{"value": 12, "label": "ペットツーリズム", "label_en": "Pet Tourism"},
	map[string]any{"value": 13, "label": "ウェルネスツーリズム", "label_en": "Wellness Tourism"},
	map[string]any{
		"value":    14,
		"label":    "人材育成・採用",
		"label_en": "Human Resource Development / Recruitment",
	},
	map[string]any{
		"value":    15,
		"label":    "AI/TECH/DX INNOVATION ZONE/ ホテルシステム・管理ツール・店舗管理",
		"label_en": "Hotel Systems / Management Tools / Store Management",
	},
	map[string]any{
		"value":    16,
		"label":    "AI/TECH/DX INNOVATION ZONE/ マーケティング・集客支援",
		"label_en": "Marketing / Customer Attraction Support",
	},
	map[string]any{
		"value":    17,
		"label":    "AI/TECH/DX INNOVATION ZONE/ データ分析・予測ツール",
		"label_en": "Data Analysis / Prediction tools",
	},
	map[string]any{
		"value":    18,
		"label":    "AI/TECH/DX INNOVATION ZONE/ サービスロボット",
		"label_en": "Service Robots",
	},
	map[string]any{
		"value":    19,
		"label":    "AI/TECH/DX INNOVATION ZONE/ AIソリューション",
		"label_en": "AI Solution",
	},
	map[string]any{
		"value":    20,
		"label":    "AI/TECH/DX INNOVATION ZONE/ 人事・労務管理システム",
		"label_en": "Personnel / Labor Management System",
	},
}

var TOPIC_SURVEY_PURCHASE_OPTIONS_FOODEX = []map[string]any{
	map[string]any{"value": 1, "label": "自分で決定する", "label_en": "Make decisions independently"},
	map[string]any{
		"value":    2,
		"label":    "決定する際に中心的な役割を果たす",
		"label_en": "Play a central role in decision-making",
	},
	map[string]any{
		"value":    3,
		"label":    "決定する際に検討メンバーになる",
		"label_en": "Be involved as a member in the decision-making process",
	},
	map[string]any{
		"value":    4,
		"label":    "関与していない",
		"label_en": "Not involved in the decision-making process",
	},
}

var TOPIC_SURVEY_PURCHASE_OPTIONS_HCJ = []map[string]any{
	map[string]any{"value": 1, "label": "自分で決定する", "label_en": "Make decisions independently"},
	map[string]any{
		"value":    2,
		"label":    "決定する際に中心的な役割を果たす",
		"label_en": "Play a central role in decision-making",
	},
	map[string]any{
		"value":    3,
		"label":    "決定する際に検討メンバーになる",
		"label_en": "Be involved as a member in the decision-making process",
	},
	map[string]any{
		"value":    4,
		"label":    "関与していない",
		"label_en": "Not involved in the decision-making process",
	},
}

var TOPIC_SURVEY_VISIT_PURPOSE_FOODEX = []map[string]any{
	map[string]any{
		"value":    1,
		"label":    "国内商品仕入れ商談のため",
		"label_en": "For domestic product sourcing and business meetings",
	},
	map[string]any{
		"value":    2,
		"label":    "海外商品仕入れ商談のため",
		"label_en": "For overseas product sourcing and business meetings",
	},
	map[string]any{
		"value":    3,
		"label":    "新製品情報を知るため",
		"label_en": "To obtain information about new products",
	},
	map[string]any{
		"value":    4,
		"label":    "新しい提携先取引先を探すため",
		"label_en": "To find new business partners and clients",
	},
	map[string]any{
		"value":    5,
		"label":    "既存顧客とのコミュニケーション",
		"label_en": "To communicate with existing customers",
	},
	map[string]any{
		"value":    6,
		"label":    "次回出展検討のため",
		"label_en": "To consider participation in the next exhibition",
	},
}

var TOPIC_SURVEY_VISIT_PURPOSE_HCJ = []map[string]any{
	map[string]any{
		"value":    1,
		"label":    "国内商品仕入れ・商談のため",
		"label_en": "For domestic product sourcing and business meetings",
	},
	map[string]any{
		"value":    2,
		"label":    "海外商品仕入れ・商談のため",
		"label_en": "For overseas product sourcing and business meetings",
	},
	map[string]any{
		"value":    3,
		"label":    "新製品情報を知るため",
		"label_en": "To obtain information about new products",
	},
	map[string]any{
		"value":    4,
		"label":    "新しい提携先・取引先を探すため",
		"label_en": "To find new business partners and clients",
	},
	map[string]any{
		"value":    5,
		"label":    "既存顧客とのコミュニケーション",
		"label_en": "To communicate with existing customers",
	},
	map[string]any{
		"value":    6,
		"label":    "次回出展検討のため",
		"label_en": "To consider participation in the next exhibition",
	},
	map[string]any{"value": 7, "label": "その他", "label_en": "Other"},
}

var TOPIC_SURVEY_BUDGET_OPTIONS_FOODEX = []map[string]any{
	map[string]any{"value": 1, "label": "100万円以下", "label_en": "Under 1 million yen"},
	map[string]any{"value": 2, "label": "500万円以下", "label_en": "Under 5 million yen"},
	map[string]any{"value": 3, "label": "1,000万円以下", "label_en": "Under 10 million yen"},
	map[string]any{"value": 4, "label": "3,000万円以下", "label_en": "Under 30 million yen"},
	map[string]any{"value": 5, "label": "5,000万円以下", "label_en": "Under 50 million yen"},
	map[string]any{"value": 6, "label": "5,001万円以上", "label_en": "Over 50.01 million yen"},
}

var TOPIC_SURVEY_BUDGET_OPTIONS_HCJ = []map[string]any{
	map[string]any{"value": 1, "label": "100万円以下", "label_en": "Under 1 million yen"},
	map[string]any{"value": 2, "label": "500万円以下", "label_en": "Under 5 million yen"},
	map[string]any{"value": 3, "label": "1,000万円以下", "label_en": "Under 10 million yen"},
	map[string]any{"value": 4, "label": "3,000万円以下", "label_en": "Under 30 million yen"},
	map[string]any{"value": 5, "label": "5,000万円以下", "label_en": "Under 50 million yen"},
	map[string]any{"value": 6, "label": "5,000万円超", "label_en": "Over 50.01 million yen"},
}
