SELECT
    name AS spot_name,
    SUBSTRING(website, '://([^/]+)') AS domain,
    COUNT(*) AS domain_count
FROM
    "MY_TABLE"
GROUP BY
    spot_name,
    SUBSTRING(website, '://([^/]+)')
HAVING
    COUNT(SUBSTRING(website, '://([^/]+)')) > 1;
-- aditionally we can add a CTE to avoid the repetition of the SUBSTRING

-- ://([^/]+) is a regular expression that matches the domain name
--- from jenkins github plugin
-- https://github.com/jenkinsci/github-plugin/blob/master/src/main/java/com/cloudbees/jenkins/GitHubRepositoryName.java#L50