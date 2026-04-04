-- Create "cps" table
CREATE TABLE `cps` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NOT NULL,
  `category` text NOT NULL,
  `link` text NULL,
  `user_cps` integer NOT NULL,
  CONSTRAINT `cps_users_cps` FOREIGN KEY (`user_cps`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "cps_name_key" to table: "cps"
CREATE UNIQUE INDEX `cps_name_key` ON `cps` (`name`);
-- Create "comments" table
CREATE TABLE `comments` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `content` text NOT NULL,
  `created_at` datetime NOT NULL,
  `cp_comments` integer NOT NULL,
  `comment_children` integer NULL,
  `user_comments` integer NOT NULL,
  CONSTRAINT `comments_users_comments` FOREIGN KEY (`user_comments`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `comments_comments_children` FOREIGN KEY (`comment_children`) REFERENCES `comments` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT `comments_cps_comments` FOREIGN KEY (`cp_comments`) REFERENCES `cps` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "tags" table
CREATE TABLE `tags` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NOT NULL,
  `user_tags` integer NOT NULL,
  CONSTRAINT `tags_users_tags` FOREIGN KEY (`user_tags`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "tags_name_key" to table: "tags"
CREATE UNIQUE INDEX `tags_name_key` ON `tags` (`name`);
-- Create "users" table
CREATE TABLE `users` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `sub` text NOT NULL
);
-- Create index "users_sub_key" to table: "users"
CREATE UNIQUE INDEX `users_sub_key` ON `users` (`sub`);
-- Create "cp_tags" table
CREATE TABLE `cp_tags` (
  `cp_id` integer NOT NULL,
  `tag_id` integer NOT NULL,
  PRIMARY KEY (`cp_id`, `tag_id`),
  CONSTRAINT `cp_tags_tag_id` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `cp_tags_cp_id` FOREIGN KEY (`cp_id`) REFERENCES `cps` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "user_liked_cps" table
CREATE TABLE `user_liked_cps` (
  `user_id` integer NOT NULL,
  `cp_id` integer NOT NULL,
  PRIMARY KEY (`user_id`, `cp_id`),
  CONSTRAINT `user_liked_cps_cp_id` FOREIGN KEY (`cp_id`) REFERENCES `cps` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT `user_liked_cps_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
);
