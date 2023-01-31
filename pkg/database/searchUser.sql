create function searchuser(search character varying, teamid bigint)
    returns TABLE(id bigint, name text, profile_picture text, isinvited boolean)
    language plpgsql
as
$$
BEGIN
    RETURN QUERY
        SELECT u.id, concat(u.first_name, ' ', u.last_name) AS name, s.profile_picture, inv.id IS NOT NULL AS isInvited
        FROM users u
                 JOIN info i ON u.info_id = i.id
                 JOIN socials s ON i.socials_id = s.id
                 LEFT JOIN invite inv ON inv.user_id = u.id AND inv.team_id = teamID
        WHERE concat(u.first_name, ' ', u.last_name) ILIKE '%' || search || '%' AND u.team_id IS NULL;
END;
$$;

alter function searchuser(varchar, bigint) owner to doadmin;

