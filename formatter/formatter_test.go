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
			"SELECT  ",
		},
		{
			"selecT",
			"SELECT  ",
		},
		{
			"SELECT 6",
			"SELECT  6",
		},
		{
			"SELECT 6, 7.89, .2",
			`SELECT  6,
        7.89,
        .2`,
		},
		{
			"SELECT (), (3), (3 4), ((1,2),3), (3,8), (x), (y,",
			`SELECT  (),
        (3),
        (3, 4),
        ((1, 2), 3),
        (3, 8),
        (x),
        (y)`,
		},
		{
			"SELeCT x,Y,z",
			`SELECT  x,
        y,
        z`,
		},
		{
			"select id.x\n,xyz.y,\nz,     BLAH., .fuck",
			`SELECT  id.x,
        xyz.y,
        z,
        blah.,
        .fuck`,
		},
		{
			"SELECT x\n AS what,y,\nkay.the AS YaY",
			`SELECT  x                                                                                                               AS what,
        y,
        kay.the                                                                                                         AS yay`,
		},
		{
			"SELECT  x  :: VARCHAR(69420) AS what,\n y :: INT, kay.the :: \nWASSUP(12,18)  AS yay, wot ::",
			`SELECT  x                                                                                             :: VARCHAR(69420) AS what,
        y                                                                                             :: INT,
        kay.the                                                                                       :: WASSUP(12, 18) AS yay,
        wot`,
		},
		{
			"SELECT  test(), COALESCE(a.name,Unknown,(4)) :: VARCHAR(300) AS area_name, trim(area_name) :: VARCHAR(150)   AS area_name, test.f_sql_GENERATE_key_from_string(area_name) AS area_key",
			`SELECT  TEST(),
        COALESCE(a.name, unknown, (4))                                                                :: VARCHAR(300)   AS area_name,
        TRIM(area_name)                                                                               :: VARCHAR(150)   AS area_name,
        test.f_sql_generate_key_from_string(area_name)                                                                  AS area_key`,
		},
	}

	for _, test := range tests {
		f := Create(test.input)
		actual := f.Format()

		if test.output != actual {
			t.Fatalf("error:\n< EXPECTED >\n%s\n\n< ACTUAL >\n%s", test.output, actual)
		} else {
			fmt.Printf("< LOGGING >\n%s\n\n", actual)
		}

	}
}

func TestGeneral(t *testing.T) {
	input := "SELECT x.what('string')"
	f := Create(input)
	t.Log(f.tokens)
}
