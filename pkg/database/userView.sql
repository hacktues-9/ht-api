create function userview(uid bigint)
    returns TABLE(first_name text, last_name text, email text, elsys_email text, mobile text, sclass text, shirt_size text, email_verified boolean, profile_pic_verified boolean, discord text, github text, looking_for_team boolean, profile_picture text)
    language plpgsql
as
$$
BEGIN
    RETURN QUERY
        SELECT users.first_name, users.last_name, users.email, users.elsys_email, users.mobile, CONCAT(info.grade, class.name) AS sclass, shirt_size.name, security.email_verified, security.manual_verified, concat(discord.username, discord.discriminator) AS discord, github.login, users.looking_for_team, socials.profile_picture
        FROM users
                 JOIN info ON users.info_id = info.id
                 LEFT JOIN security ON users.security_id = security.id
                 LEFT JOIN shirt_size ON info.shirt_size_id = shirt_size.id
                 JOIN class ON info.class_id = class.id
                 JOIN socials on info.socials_id = socials.id
                 LEFT JOIN discord ON socials.discord_id = discord.id
                 LEFT JOIN github ON socials.github_id = github.id
        WHERE users.id = uID;
END;
$$;

alter function userview(bigint) owner to postgres;

