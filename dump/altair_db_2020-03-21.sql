-- phpMyAdmin SQL Dump
-- version 4.8.5
-- https://www.phpmyadmin.net/
--
-- Хост: localhost
-- Время создания: Мар 21 2020 г., 20:38
-- Версия сервера: 8.0.13-4
-- Версия PHP: 7.2.24-0ubuntu0.18.04.3

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- База данных: `s5JWeA8krh`
--
CREATE DATABASE IF NOT EXISTS `s5JWeA8krh` DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci;
USE `s5JWeA8krh`;

-- --------------------------------------------------------

--
-- Структура таблицы `ads`
--

CREATE TABLE `ads` (
  `ad_id` int(10) UNSIGNED NOT NULL,
  `title` varchar(255) NOT NULL,
  `slug` varchar(255) NOT NULL,
  `cat_id` int(10) UNSIGNED NOT NULL,
  `user_id` int(10) UNSIGNED NOT NULL DEFAULT '0',
  `text` text CHARACTER SET utf8 COLLATE utf8_general_ci,
  `price` int(10) UNSIGNED NOT NULL DEFAULT '0',
  `ip` varchar(255) NOT NULL DEFAULT '',
  `is_disabled` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='объявления';

--
-- Дамп данных таблицы `ads`
--

INSERT INTO `ads` (`ad_id`, `title`, `slug`, `cat_id`, `user_id`, `text`, `price`, `ip`, `is_disabled`, `created_at`, `updated_at`) VALUES
(8, 'мой заголовок2', 'moi-zagolovok2_8', 59, 1, 'test', 10, '', 1, '2020-03-11 00:49:40', '2020-03-19 02:31:54'),
(9, 'мой заголовок2', 'moi-zagolovok2_9', 1, 0, 'wewer', 0, '', 0, '2020-03-11 01:05:59', '2020-03-11 01:05:59'),
(10, 'мой заголовок3', 'moi-zagolovok3_10', 1, 0, '', 0, '', 0, '2020-03-11 09:44:47', '2020-03-11 09:44:47'),
(11, 'мой заголовок4', 'moi-zagolovok4_11', 1, 0, '', 0, '', 0, '2020-03-11 09:52:58', '2020-03-11 09:52:59'),
(12, 'мой заголовок4', 'moi-zagolovok4_12', 1, 0, '', 0, '', 0, '2020-03-11 09:54:08', '2020-03-11 09:54:08'),
(13, 'мой заголовок5', 'moi-zagolovok5_13', 1, 0, '', 0, '', 0, '2020-03-11 09:56:41', '2020-03-11 09:56:41'),
(14, 'mytest1', 'mytest1_14', 1, 0, '', 0, '', 0, '2020-03-11 22:52:26', '2020-03-11 22:52:27'),
(15, 'vjkur', 'vkr_15', 1, 0, '', 0, '', 0, '2020-03-14 14:16:16', '2020-03-14 14:16:16'),
(18, 'amguo', 'amgo_18', 1, 0, '', 0, '', 0, '2020-03-14 15:10:43', '2020-03-14 15:10:43'),
(21, 'qhcyn', 'qhcyn_21', 1, 0, '', 0, '', 0, '2020-03-14 15:15:49', '2020-03-14 15:15:50');

-- --------------------------------------------------------

--
-- Структура таблицы `cats`
--

CREATE TABLE `cats` (
  `cat_id` int(10) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `slug` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `parent_id` int(10) UNSIGNED NOT NULL DEFAULT '0',
  `pos` tinyint(3) UNSIGNED NOT NULL DEFAULT '1',
  `is_disabled` tinyint(4) NOT NULL DEFAULT '1'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='категории';

--
-- Дамп данных таблицы `cats`
--

INSERT INTO `cats` (`cat_id`, `name`, `slug`, `parent_id`, `pos`, `is_disabled`) VALUES
(51, 'Категория1', 'kategoriia1', 0, 1, 0),
(52, 'Категория2', 'kategoriia2', 0, 2, 0),
(53, 'Категория3', 'kategoriia3', 0, 3, 1),
(54, 'Категория1-1', 'kategoriia1-1', 51, 1, 0),
(55, 'Категория1-2', 'kategoriia1-2', 51, 2, 0),
(56, 'Категория1-3', 'kategoriia1-3', 51, 3, 0),
(57, 'Категория1-1-1', 'kategoriia1-1-1', 54, 1, 0),
(58, 'Категория1-1-2', 'kategoriia1-1-2', 54, 2, 0),
(59, 'Категория1-1-3', 'kategoriia1-1-3', 54, 3, 0),
(60, 'Категория1-2-1', 'kategoriia1-2-1', 55, 1, 0),
(61, 'Категория2-1', 'kategoriia2-1', 52, 1, 0),
(62, 'Категория2-2', 'kategoriia2-2', 52, 2, 0),
(63, 'Категория2-2-1', 'kategoriia2-2-1', 62, 1, 0);

-- --------------------------------------------------------

--
-- Структура таблицы `cats_properties`
--

CREATE TABLE `cats_properties` (
  `cat_id` int(11) NOT NULL,
  `property_id` int(11) NOT NULL,
  `pos` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='связь категории и его аттрибутов';

--
-- Дамп данных таблицы `cats_properties`
--

INSERT INTO `cats_properties` (`cat_id`, `property_id`, `pos`) VALUES
(53, 28, 2);

-- --------------------------------------------------------

--
-- Структура таблицы `images`
--

CREATE TABLE `images` (
  `img_id` int(10) UNSIGNED NOT NULL,
  `filepath` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `el_id` int(10) UNSIGNED NOT NULL,
  `is_disabled` tinyint(1) NOT NULL DEFAULT '1',
  `opt` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='общая таблица для всех картинок';

--
-- Дамп данных таблицы `images`
--

INSERT INTO `images` (`img_id`, `filepath`, `el_id`, `is_disabled`, `opt`, `created_at`) VALUES
(46, '82/8e/a6059bc2803959a8df9e9a168a711503fa331a3dd75279d5bf5fe8e3ee2e_vykb.jpg', 8, 0, 'ad', '2020-03-15 23:32:59'),
(47, '2c/27/976179626a8cb34bcec549930b042100054c212edbd0beb4f2ed4c043641_wkhs.jpg', 8, 0, 'ad', '2020-03-15 23:32:59'),
(48, '2d/d5/d6ea5ae8c87d7cf75ad28743fb4836e2a862781434757d3d8deb5beb746b_hweh.jpg', 8, 0, 'ad', '2020-03-15 23:33:01');

-- --------------------------------------------------------

--
-- Структура таблицы `kind_properties`
--

CREATE TABLE `kind_properties` (
  `kind_property_id` int(10) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='виды тегов свойств к категориям';

--
-- Дамп данных таблицы `kind_properties`
--

INSERT INTO `kind_properties` (`kind_property_id`, `name`) VALUES
(4, 'checkbox'),
(7, 'input_date'),
(8, 'input_datetime'),
(2, 'input_number'),
(9, 'photo'),
(5, 'radio'),
(3, 'select'),
(6, 'textarea');

-- --------------------------------------------------------

--
-- Структура таблицы `properties`
--

CREATE TABLE `properties` (
  `property_id` int(10) UNSIGNED NOT NULL,
  `title` varchar(255) NOT NULL,
  `kind_property_id` int(10) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `is_require` tinyint(1) NOT NULL DEFAULT '0',
  `is_can_as_filter` tinyint(1) NOT NULL DEFAULT '0',
  `max_int` int(11) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='свойства категории';

--
-- Дамп данных таблицы `properties`
--

INSERT INTO `properties` (`property_id`, `title`, `kind_property_id`, `name`, `is_require`, `is_can_as_filter`, `max_int`) VALUES
(28, 'Цвет', 5, 'color', 0, 0, 0);

-- --------------------------------------------------------

--
-- Структура таблицы `users`
--

CREATE TABLE `users` (
  `user_id` int(10) UNSIGNED NOT NULL,
  `email` varchar(255) NOT NULL,
  `email_is_confirmed` tinyint(1) NOT NULL DEFAULT '0',
  `name` varchar(255) NOT NULL DEFAULT '',
  `password` varchar(255) NOT NULL,
  `avatar` varchar(255) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Дамп данных таблицы `users`
--

INSERT INTO `users` (`user_id`, `email`, `email_is_confirmed`, `name`, `password`, `avatar`, `created_at`, `updated_at`) VALUES
(33, 'test@rkwxk.thm', 1, 'Lorem3', '$2a$04$CE/EzfZNijUS7IjG2buOW.PLT7tgXAJdLQFmboQcxTGFnGR.Bw/Pa', 'ea/d4/5aebfa476563cdc0772f6ad978fde3d8301f461ef308854c3c5d939635d5_pyol.jpg', '2020-03-09 01:02:04', '2020-03-15 16:08:38'),
(35, 'test@kotud.myq', 0, '', '$2a$04$ynFohhuZJRR7cXPxtXmBVe8UAYG6cwKp88sPxXtCcd5Jnx1/B.yqe', '', '2020-03-09 01:05:32', '2020-03-09 01:05:32');

-- --------------------------------------------------------

--
-- Структура таблицы `value_properties`
--

CREATE TABLE `value_properties` (
  `value_id` int(11) NOT NULL,
  `title` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `pos` int(10) UNSIGNED NOT NULL,
  `property_id` int(10) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='значения для свойств(а)' ROW_FORMAT=COMPACT;

--
-- Дамп данных таблицы `value_properties`
--

INSERT INTO `value_properties` (`value_id`, `title`, `pos`, `property_id`) VALUES
(19, 'test1', 1, 28),
(20, 'test2', 2, 28),
(21, 'test3', 3, 28);

--
-- Индексы сохранённых таблиц
--

--
-- Индексы таблицы `ads`
--
ALTER TABLE `ads`
  ADD PRIMARY KEY (`ad_id`),
  ADD UNIQUE KEY `slug` (`slug`);

--
-- Индексы таблицы `cats`
--
ALTER TABLE `cats`
  ADD PRIMARY KEY (`cat_id`),
  ADD UNIQUE KEY `parent_id` (`parent_id`,`slug`);

--
-- Индексы таблицы `images`
--
ALTER TABLE `images`
  ADD PRIMARY KEY (`img_id`);

--
-- Индексы таблицы `kind_properties`
--
ALTER TABLE `kind_properties`
  ADD PRIMARY KEY (`kind_property_id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- Индексы таблицы `properties`
--
ALTER TABLE `properties`
  ADD PRIMARY KEY (`property_id`);

--
-- Индексы таблицы `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`user_id`),
  ADD UNIQUE KEY `email` (`email`);

--
-- Индексы таблицы `value_properties`
--
ALTER TABLE `value_properties`
  ADD PRIMARY KEY (`value_id`);

--
-- AUTO_INCREMENT для сохранённых таблиц
--

--
-- AUTO_INCREMENT для таблицы `ads`
--
ALTER TABLE `ads`
  MODIFY `ad_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=77;

--
-- AUTO_INCREMENT для таблицы `cats`
--
ALTER TABLE `cats`
  MODIFY `cat_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=112;

--
-- AUTO_INCREMENT для таблицы `images`
--
ALTER TABLE `images`
  MODIFY `img_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=49;

--
-- AUTO_INCREMENT для таблицы `kind_properties`
--
ALTER TABLE `kind_properties`
  MODIFY `kind_property_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=50;

--
-- AUTO_INCREMENT для таблицы `properties`
--
ALTER TABLE `properties`
  MODIFY `property_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=31;

--
-- AUTO_INCREMENT для таблицы `users`
--
ALTER TABLE `users`
  MODIFY `user_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=84;

--
-- AUTO_INCREMENT для таблицы `value_properties`
--
ALTER TABLE `value_properties`
  MODIFY `value_id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=22;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
