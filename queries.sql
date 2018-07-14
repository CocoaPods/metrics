SELECT dependency_name, dependency_version, to_char(sent_at, 'YYYY-MM') AS rollup_date,
       count(CASE WHEN pod_try THEN 1 ELSE null END) AS pod_tries,
       count(CASE WHEN pod_try THEN null ELSE 1 END) AS downloads
  FROM cocoapods_stats_production.install
  GROUP BY rollup_date, dependency_name, dependency_version;


SELECT dependency_name, dependency_version, product_type, COUNT(DISTINCT(user_id)) AS installs, to_char(sent_at, 'YYYY-MM') AS rollup_date
  FROM cocoapods_stats_production.install
 WHERE pod_try = false
GROUP BY rollup_date, dependency_version, dependency_name, product_type;
