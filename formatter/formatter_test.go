package formatter

import (
	"fmt"
	"testing"
)

func TestBasicSelect(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{
			"SELECT",
			"SELECT  \n;",
		},
		{
			"selecT",
			"SELECT  \n;",
		},
		{
			"SELECT 6",
			"SELECT  6\n;",
		},
		{
			"SELECT 10 SELECT 'what'",
			"SELECT  10\n;\nSELECT  'what'\n;",
		},
		{
			"SELECT 10; SELECT 'what'\n;",
			"SELECT  10\n;\nSELECT  'what'\n;",
		},
		{
			"SELECT 6, 7.89, .2",
			`SELECT  6,
        7.89,
        .2
;`,
		},
		{
			"SELECT (), (3), (3 4), ((1,2),3), (3,8), (x), (y,",
			`SELECT  (),
        (3),
        (3 4),
        ((1, 2), 3),
        (3, 8),
        (x),
        (y, )
;`,
		},
		{
			"SELeCT x,Y,z",
			`SELECT  x,
        y,
        z
;`,
		},
		{
			"select id.x\n,xyz.y,\nz,     BLAH., .fuck",
			`SELECT  id.x,
        xyz.y,
        z,
        blah.,
        .fuck
;`,
		},
		{
			"SELECT x\n AS what,y,\nkay.the AS YaY",
			`SELECT  x                                                                                                               AS what,
        y,
        kay.the                                                                                                         AS yay
;`,
		},
		{
			"SELECT  x  :: VARCHAR(69420) AS what,\n y :: INT, kay.the :: \nWASSUP(12,18)  AS yay, wot ::",
			`SELECT  x                                                                                             :: VARCHAR(69420) AS what,
        y :: INT,
        kay.the                                                                                       :: WASSUP(12, 18) AS yay,
        wot ::
;`,
		},
		{
			"SELECT  test(), COALESCE(a.name,'Unknown',(4)) :: VARCHAR(300) AS area_name, trim(area_name) :: VARCHAR(150)   AS area_name, test.f_sql_GENERATE_key_from_string(area_name) AS area_key",
			`SELECT  TEST(),
        COALESCE(a.name, 'Unknown', (4))                                                              :: VARCHAR(300)   AS area_name,
        TRIM(area_name)                                                                               :: VARCHAR(150)   AS area_name,
        test.f_sql_generate_key_from_string(area_name)                                                                  AS area_key
;`,
		},
		{
			"SELECT  (city_sk|| '|'|| odl.f_sql_gnerate_nk_from_StrinG(area_name) || '|' ||area_nk ) ::   VARCHAR(300)  AS area_sk, area_nk",
			`SELECT  (city_sk || '|' || odl.f_sql_gnerate_nk_from_string(area_name) || '|' || area_nk)             :: VARCHAR(300)   AS area_sk,
        area_nk
;`,
		},
		{
			`SELECT  (city_sk || '|' || odl.f_sql_gnerate_nk_from_string(area_name) || '|' || area_nk)      :: VARCHAR(300)  AS area_sk,
        area_nk,
        tenant_sk,
        city_sk,
        zone_sk 
        etl_source_sk,
        trim(area_name)                                                                                :: VARCHAR(150)  AS area_name,
        TRIM(area_name_ar)                                                                :: VARCHAR(150)  AS area_name_ar,
        area_reference_number
SELECT      'unknown'                                                                                      AS area_sk,
        0                                                                                       AS area_nk,
        'unknown'                                                                                      AS tenant_sk,
        'unknown'                                                                                      AS city_sk,
        'unknown'                                                                                      AS zone_sk,
        'unknown'                                                                    AS etl_source_sk,
        'Unknown'                                                                                      AS area_name,
        'Unknown'                                                                         AS area_name_ar,
        0                                                                                              AS area_reference_number`,
			`SELECT  (city_sk || '|' || odl.f_sql_gnerate_nk_from_string(area_name) || '|' || area_nk)             :: VARCHAR(300)   AS area_sk,
        area_nk,
        tenant_sk,
        city_sk,
        zone_sk etl_source_sk,
        TRIM(area_name)                                                                               :: VARCHAR(150)   AS area_name,
        TRIM(area_name_ar)                                                                            :: VARCHAR(150)   AS area_name_ar,
        area_reference_number
;
SELECT  'unknown'                                                                                                       AS area_sk,
        0                                                                                                               AS area_nk,
        'unknown'                                                                                                       AS tenant_sk,
        'unknown'                                                                                                       AS city_sk,
        'unknown'                                                                                                       AS zone_sk,
        'unknown'                                                                                                       AS etl_source_sk,
        'Unknown'                                                                                                       AS area_name,
        'Unknown'                                                                                                       AS area_name_ar,
        0                                                                                                               AS area_reference_number
;`,
		},
		{
			"SELECT \"dev_modl\".f_sql_gnerate_nk_from_string(t.trim) AS \"car_trim_new_nk\",",
			"SELECT  \"dev_modl\".f_sql_gnerate_nk_from_string(t.trim)                                                                 AS \"car_trim_new_nk\"\n;",
		},
		{
			"SELECT  FALSE, TRUe, null",
			`SELECT  FALSE,
        TRUE,
        NULL
;`,
		},
		{
			"SELECT  ROW_NUMBER() OVER (PARTITION by row1 ORdER BY xyz DESC)",
			`SELECT  ROW_NUMBER() OVER (PARTITION BY row1 ORDER BY xyz DESC)
;`,
		},
		{
			"SELECT LAG(date_ending_nk, 1)\n OVER (PARTITION BY c.user_sk ORDER BY date_sign_nk, contract_nk )   :: DATE            AS previous_contract_date_ending,",
			`SELECT  LAG(date_ending_nk, 1) OVER (PARTITION BY c.user_sk ORDER BY date_sign_nk, contract_nk)       :: DATE           AS previous_contract_date_ending
;`,
		},
		{
			`SELECT  x, -- comment
        xyz + -- comment
        4,
        (what - the || --okay
        'string'),
        bla :: INT AS -- this is bad
        col,
        -- line column
        another_line`,
			`SELECT  x, -- comment
        xyz + 4, -- comment
        (what - the || 'string'), --okay
        bla                                                                                           :: INT            AS col, -- line column
        another_line
;`,
		},
		{
			"SELECT  -xyz + -  asdf (-xyz || asdf - 'asdf') + (- *)+ (-'asdf') - :: INT AS what,",
			`SELECT  -xyz + -asdf (-xyz || asdf - 'asdf') + ( -  * ) + ( - 'asdf') -                               :: INT            AS what
;`,
		},
	}

	for i := 0; i < len(tests); i++ {
		f := Create(tests[i].input)
		actual := f.Format()

		if tests[i].output != actual {
			fmt.Printf("error [%d]:\n< EXPECTED >\n%s\n\n< ACTUAL >\n%s", i, tests[i].output, actual)
			fmt.Printf("\n< EXPECTED BYTES >\n%v\n\n< ACTUAL BYTES >\n%v", []byte(tests[i].output), []byte(actual))
			t.Fatal()
		} else {
			fmt.Printf("< LOGGING [%d] >\n%s\n\n", i, actual)
		}

	}
}

func TestGeneral(t *testing.T) {
	input := `
SELECT -xyz + -  asdf (-xyz || asdf - 'asdf') + (- *)+ (-'asdf') - ++asdf :: INT AS what,
`
	f := Create(input)
	fmt.Println(f.tokens)
	formatted := f.Format()
	fmt.Println(formatted)
}
