create table simple_rest_app.user_role (
	rec_dttm	timestamp default now()::timestamp,
	id			serial primary key,
	code 		int2 unique not null,
	name 		text not null
);

comment on table simple_rest_app.user_role is 'Список доступных ролей для пользователей.';

comment on column simple_rest_app.user_role.rec_dttm is 'Дата и время записи сообщения.';
comment on column simple_rest_app.user_role.id is 'Идентификатор роли пользователя.';
comment on column simple_rest_app.user_role.code is 'Код роли пользователя.';
comment on column simple_rest_app.user_role.name is 'Наименование роли.';

insert into simple_rest_app.user_role(
	code,
	name
)
	select
		1,
		'Администратор';
		
insert into simple_rest_app.user_role(
	code,
	name
)
	select
		0,
		'Пользователь';

	create table simple_rest_app.user (
	rec_dttm	timestamp default now()::timestamp,
	id			serial primary key,
	code_role	int2 default 0,
	name		text unique not null,
	password	text not null,
	foreign key (code_role) references simple_rest_app.user_role (code) on delete no action on update cascade
);

comment on table simple_rest_app.user is 'Перечень пользователей приложения.';

comment on column simple_rest_app.user.rec_dttm is 'Дата и время записи сообщения.';
comment on column simple_rest_app.user.id is 'Идентификатор пользователея.';
comment on column simple_rest_app.user.name is 'Логин пользователя.';
comment on column simple_rest_app.user.password is 'Пароль пользователя.';
comment on column simple_rest_app.user.code_role is 'Код роли польщователя.';

insert into simple_rest_app.user(
	name,
	password,
	code_role
)
	select
		'admin',
		'admin',
		1;
	
		
create table simple_rest_app.post (
	rec_dttm timestamp default now()::timestamp,
	id serial primary key,
	id_user int4,
	message text not null,
	foreign key (id_user) references simple_rest_app.user (id) on delete cascade on update cascade
)

comment on table simple_rest_app.post is 'Перечень постов приложения.';

comment on column simple_rest_app.post.rec_dttm is 'Дата и время записи сообщения.';
comment on column simple_rest_app.post.id is 'Идентификатор поста.';
comment on column simple_rest_app.post.id_user is 'Идентификатор пользователя, создавшего пост.';
comment on column simple_rest_app.post.message is 'Содержание поста.';

create index idxpstidusr on simple_rest_app.post using btree (id_user);

create table simple_rest_app.post_comment(
	rec_dttm	timestamp default now()::timestamp,
	id			serial primary key,
	id_user		int4,
	id_post		int4,
	comment		text not null,
	foreign key (id_user) references simple_rest_app.user (id) on delete cascade on update cascade,
	foreign key (id_post) references simple_rest_app.post (id) on delete cascade on update cascade
);

comment on table simple_rest_app.post_comment is 'Таблица, хранящая комментарии на посты пользователей.';

comment on column simple_rest_app.post_comment.rec_dttm is 'Дата и время записи сообщения.';
comment on column simple_rest_app.post_comment.id is 'Идентификатор комментария.';
comment on column simple_rest_app.post_comment.id_user is 'Внешний ключ, ссылающийся на таблицу simple_rest_app.user, идентификатор пользователя, оставившего комментарий.';
comment on column simple_rest_app.post_comment.id_post is 'Внешний ключ, ссылающийся на таблицу simple_rest_app.post, идентификатор поста.';
comment on column simple_rest_app.post_comment.comment is 'Текст комментария.';

create index idxpstcommidusr on simple_rest_app.post_comment using btree (id_user);
create index idxpstcommidpst on simple_rest_app.post_comment using btree (id_post);