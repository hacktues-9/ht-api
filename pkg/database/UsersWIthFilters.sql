 -- PostgreSQL
 -- check if any of the parameters are null or empty
    -- if they are, then don't filter on that column
    -- if they aren't, then filter on that column
    -- name, email, mobile, elsys_email check on like instead of = because they are text
    -- name is combined with first and last name
create function get_user_filtered( Ishirt_size text, Igrade bigint, Iclass text, Iname text, Iemail text, Imobile text, Ielsys_email text, IteamName text, IeatingPreference text)
returns TABLE(class text, first_name text, last_name text, email text, elsys_email text, mobile text, shirt_size text, eating_preference text, email_verified text, elsys_email_verified text, manual_verified text, discord text, github text, team text)
as $$
begin
    return query
        select concat(i.grade, c.name) as class, u.first_name as first_name, u.last_name as last_name, u.email as email, u.elsys_email as elsys_email, u.mobile as mobile, ss.name as shirt_size, ep.name as eating_preference, s.email_verified as email_verified, s.elsys_email_verified as elsys_email_verified, s.manual_verified as manual_verified, concat(d.username,'#',d.discriminator) as discord, g.name as github, t.name as team
        from users u
        JOIN info i on i.id = u.info_id
        JOIN shirt_size ss on i.shirt_size_id = ss.id
        JOIN class c on c.id = i.class_id
        JOIN team t on t.id = u.team_id
        JOIN eating_preference ep on ep.id = i.eating_preference_id
        JOIN security s on u.security_id = s.id
        JOIN socials s2 on i.socials_id = s2.id
        JOIN discord d on s2.discord_id = d.id
        JOIN github g on s2.github_id = g.id
        where
        (Ishirt_size is null or Ishirt_size = '' or ss.name = Ishirt_size)
        and (Igrade is null or Igrade = 0 or i.grade = Igrade)
        and (Iclass is null or Iclass = '' or c.name = Iclass)
        and (IeatingPreference is null or IeatingPreference = '' or ep.name = IeatingPreference)
        and (Iname is null or Iname = '' or (u.first_name || ' ' || u.last_name) like '%' || Iname || '%')
        and (Iemail is null or Iemail = '' or u.email like '%' || Iemail || '%')
        and (Imobile is null or Imobile = '' or u.mobile like '%' || Imobile || '%')
        and (Ielsys_email is null or Ielsys_email = '' or u.elsys_email like '%' || Ielsys_email || '%')
        and (IteamName is null or IteamName = '' or t.name like '%' || IteamName || '%')
        and u.deleted_at is null;

end;
$$ language plpgsql;