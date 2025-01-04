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


SELECT  (city_sk || '|' || odl.f_sql_gnerate_nk_from_string(area_name) || '|' || area_nk)             :: VARCHAR(300)   AS area_sk,
        area_nk,
        tenant_sk,
        city_sk,
        zone_sk,
        etl_source_sk,
        TRIM(area_name)                                                                               :: VARCHAR(150)   AS area_name,
        TRIM(area_name_ar)                                                                            :: VARCHAR(150)   AS area_name_ar,
        area_reference_number
SELECT  'unknown'                                                                                                       AS area_sk,
        0              G                                                                                                 AS area_nk,
        'unknown'                                                                                                       AS tenant_sk,
        'unknown'                                                                                                       AS city_sk,
        'unknown'                                                                                                       AS zone_sk,
        'unknown'                                                                                                       AS etl_source_sk,
        'Unknown'                                                                                                       AS area_name,
        'Unknown'                                                                                                       AS area_name_ar,
        0                                                                                                               AS area_reference_number
;

SELECT  (city_sk || '|' || odl.f_sql_gnerate_nk_from_string(area_name) || '|' || area_nk)      :: VARCHAR(300)  AS area_sk,
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
        '                 unknown'                                                                    AS etl_source_sk,
        'Unknown'                                                                                      AS area_name,
        'Unknown'                                                                         AS area_name_ar,
        0                                                                                              AS area_reference_number
;

SELECT  (LAG(date_ending_nk, 1) OVER (PARTITION BY c.user_sk ORDER BY date_sign_nk, contract_nk))     :: DATE           AS previous_contract_date_ending,

SELECT  x                                                                                             :: VARCHAR(69420) AS what,
        y :: INT,
        kay.the                                                                                       :: WASSUP(12, 18) AS yay,
        wot :: 
