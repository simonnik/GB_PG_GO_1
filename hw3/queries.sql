select * from posts
                  left join posts_images pi on posts.id = pi.post_id
where posts.is_active = true
limit 50;

select * from posts_comments where post_id = 3;

select count(*), post_id from posts_comments
where post_id in (1,2,3,4217342,6312989)
group by post_id;


select * from users where id = 3;

select * from posts_favorites where post_id=3;