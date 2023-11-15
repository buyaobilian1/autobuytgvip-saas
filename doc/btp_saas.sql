

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for order
-- ----------------------------
DROP TABLE IF EXISTS `order`;
CREATE TABLE `order`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL COMMENT '买家user id',
  `agent_user_id` int NOT NULL COMMENT '下单时机器人所属代理user id，机器人可能更换',
  `bot_id` bigint NOT NULL COMMENT '机器人id',
  `order_no` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '订单号',
  `fragment_ref_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '与fragment对应的ref',
  `receive_tg_username` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '接收会员的tg用户名',
  `vip_month` int NOT NULL COMMENT '开通vip时长3、6、12',
  `usdt_amount` decimal(10, 2) NOT NULL COMMENT '需支付USDT数',
  `base_amount` decimal(10, 2) NOT NULL COMMENT '给到代理的价格',
  `status` int NOT NULL DEFAULT 1 COMMENT '1：等待支付，2：支付成功，3：开通成功，4：过期作废',
  `agent_status` int NOT NULL DEFAULT 1 COMMENT '1：等待入账，2：代理已入账',
  `created_at` timestamp NOT NULL COMMENT '订单创建时间',
  `expired_at` timestamp NOT NULL COMMENT '订单过期时间',
  `tg_chat_id` bigint NOT NULL,
  `tg_msg_id` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `order_no_unique`(`order_no` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 76 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '购物订单' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for param
-- ----------------------------
DROP TABLE IF EXISTS `param`;
CREATE TABLE `param`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `k` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `v1` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v2` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v3` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v4` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v5` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `v6` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '系统配置' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for recharge
-- ----------------------------
DROP TABLE IF EXISTS `recharge`;
CREATE TABLE `recharge`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL COMMENT 'user id',
  `bot_id` bigint NOT NULL COMMENT '机器人id',
  `order_no` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '订单号',
  `amount` decimal(10, 2) NOT NULL COMMENT '需支付USDT数',
  `status` int NOT NULL DEFAULT 1 COMMENT '1：等待支付，2：支付成功，3：过期作废',
  `actual_amount` decimal(10, 2) NOT NULL COMMENT '实际支付金额，由支付平台返回',
  `created_at` timestamp NOT NULL COMMENT '订单创建时间',
  `tg_chat_id` bigint NOT NULL,
  `tg_msg_id` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `order_no_unique`(`order_no` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 26 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '充值记录' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `parent_id` int NOT NULL DEFAULT 0 COMMENT '上级用户id，从哪个用户的机器人进来的',
  `tg_id` bigint NOT NULL COMMENT 'tg id',
  `tg_username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'tg username',
  `balance` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT 'USDT钱包',
  `brokerage` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '代理佣金',
  `tron_addr` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'USDT收款钱包地址',
  `created_at` timestamp NOT NULL COMMENT '创建时间',
  `bot_id` bigint NULL DEFAULT NULL COMMENT '机器人id',
  `bot_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '机器人token',
  `bot_username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '机器人username',
  `bot_status` int NOT NULL DEFAULT 1 COMMENT '1：正常，2：无，3：异常，4：禁用',
  `bot_created_at` timestamp NULL DEFAULT NULL COMMENT '机器人添加时间',
  `three_month_price` decimal(10, 2) NULL DEFAULT NULL COMMENT '三个月套餐价格',
  `six_month_price` decimal(10, 2) NULL DEFAULT NULL COMMENT '六个月套餐价格',
  `twelve_month_price` decimal(10, 2) NULL DEFAULT NULL COMMENT '一年套餐价格',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for withdraw
-- ----------------------------
DROP TABLE IF EXISTS `withdraw`;
CREATE TABLE `withdraw`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `tg_id` bigint NOT NULL,
  `order_no` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '提现订单号',
  `amount` decimal(10, 2) NOT NULL COMMENT '提现金额',
  `tron_addr` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '提现钱包',
  `status` int NOT NULL COMMENT '1：等待提现，2：提现完成，3：提现失败',
  `created_at` timestamp NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '提现记录' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;


INSERT INTO `btp_saas`.`param` (`k`, `v1`, `v2`, `v3`, `v4`, `v5`, `v6`, `remark`) VALUES ('base_price', '13', '17', '30', NULL, NULL, NULL, '平台底价，v1、v2、v3分别为3,6,12个月的会员价');
