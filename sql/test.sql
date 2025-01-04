SELECT  x                                                                                                               AS what,
        y,
        kay.the                                                                                                         AS yay

SELECT  x                                                                                             :: VARCHAR(69420) AS what
        y                                                                                             :: INT,
        kay.the                                                                                       :: WASSUP(12,18)  AS yay

SELECT  COALESCE(a.name, 'Unknown')                                                                   :: VARCHAR(300)   AS area_name,
        TRIM(area_name)                                                                               :: VARCHAR(150)   AS area_name,
        test.f_sql_gnerate_key_from_string(area_name)                                                                   AS area_key

SELECT  COALESCE(a.name, 'Unknown') :: VARCHAR(300) AS area_name, TRIM(area_name) :: VARCHAR(150)   AS area_name, test.f_sql_gnerate_key_from_string(area_name) AS area_key


SELECT  (city_sk || '|' || odl.f_sql_gnerate_key_from_string(area_name) || '|' || area_nk)             :: VARCHAR(300)   AS area_sk,
        area_nk