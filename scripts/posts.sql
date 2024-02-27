--CREATE TABLE posts (
                       --id SERIAL PRIMARY KEY,
                      -- title VARCHAR(100) NOT NULL,
                    --   owner_user_id SERIAL NOT NULL,
                  --     amount integer,
                --       description text,
              --         worker_user_id integer,
            --           post_type varchar(150),
          --             status varchar(150),
        --               payment_status varchar(150),
      --                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    --                   updated_at TIMESTAMP NULL,
  --                     deleted_at TIMESTAMP NULL
--);


/*

 -- public.posts definition

-- Drop table

-- DROP TABLE public.posts;

CREATE TABLE public.posts (
	id serial4 NOT NULL,
	title varchar NOT NULL,
	owner_user_id serial4 NOT NULL,
	amount money NULL,
	description text NULL,
	worker_user_id int8 NULL,
	post_type varchar NULL,
	status varchar NULL,
	payment_status varchar NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	deleted_at timestamp NULL
);


-- public.posts foreign keys

ALTER TABLE public.posts ADD CONSTRAINT posts_fk FOREIGN KEY (owner_user_id) REFERENCES public.users(id);
ALTER TABLE public.posts ADD CONSTRAINT posts_fk_1 FOREIGN KEY (worker_user_id) REFERENCES public.users(id);

 */