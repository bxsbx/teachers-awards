CREATE TABLE "teachers_awards"."activity"
(
    "activity_id" BIGINT IDENTITY(79, 1) NOT NULL,
    "activity_name" VARCHAR(50) NOT NULL,
    "year" BIGINT NOT NULL,
    "description" VARCHAR(1500) NOT NULL,
    "url" VARCHAR(255) NOT NULL,
    "review_num" TINYINT NOT NULL,
    "start_time" TIMESTAMP(0) NOT NULL,
    "end_time" TIMESTAMP(0) NOT NULL,
    "create_time" TIMESTAMP(0) NOT NULL,
    "update_time" TIMESTAMP(0),
    "delete_time" TIMESTAMP(0),
    NOT CLUSTER PRIMARY KEY("activity_id"),
    CHECK("year" >= 0)) STORAGE(ON "teachers_awards", CLUSTERBTR) ;

COMMENT ON TABLE "teachers_awards"."activity" IS '名师评优活动';
COMMENT ON COLUMN "teachers_awards"."activity"."activity_id" IS '主键id';
COMMENT ON COLUMN "teachers_awards"."activity"."activity_name" IS '活动名称';
COMMENT ON COLUMN "teachers_awards"."activity"."year" IS '年份';
COMMENT ON COLUMN "teachers_awards"."activity"."description" IS '申报须知';
COMMENT ON COLUMN "teachers_awards"."activity"."url" IS '活动文件url';
COMMENT ON COLUMN "teachers_awards"."activity"."review_num" IS '审核人数';
COMMENT ON COLUMN "teachers_awards"."activity"."start_time" IS '开始时间';
COMMENT ON COLUMN "teachers_awards"."activity"."end_time" IS '结束时间';
COMMENT ON COLUMN "teachers_awards"."activity"."create_time" IS '创建时间';
COMMENT ON COLUMN "teachers_awards"."activity"."update_time" IS '更新时间';
COMMENT ON COLUMN "teachers_awards"."activity"."delete_time" IS '删除时间';


CREATE TABLE "teachers_awards"."activity_one_indicator"
(
    "activity_id" BIGINT NOT NULL,
    "one_indicator_id" BIGINT NOT NULL,
    "one_indicator_name" VARCHAR(50) NOT NULL,
    "content" VARCHAR(1500) NOT NULL,
    "create_time" TIMESTAMP(0) NOT NULL,
    NOT CLUSTER PRIMARY KEY("activity_id", "one_indicator_id"),
    CHECK("activity_id" >= 0)
    ,CHECK("one_indicator_id" >= 0)) STORAGE(ON "teachers_awards", CLUSTERBTR) ;

COMMENT ON TABLE "teachers_awards"."activity_one_indicator" IS '活动下的指标';
COMMENT ON COLUMN "teachers_awards"."activity_one_indicator"."activity_id" IS '活动id';
COMMENT ON COLUMN "teachers_awards"."activity_one_indicator"."one_indicator_id" IS '一级指标id';
COMMENT ON COLUMN "teachers_awards"."activity_one_indicator"."one_indicator_name" IS '一级指标名称';
COMMENT ON COLUMN "teachers_awards"."activity_one_indicator"."content" IS '评分标准说明';
COMMENT ON COLUMN "teachers_awards"."activity_one_indicator"."create_time" IS '创建时间';


CREATE TABLE "teachers_awards"."activity_two_indicator"
(
    "activity_id" BIGINT NOT NULL,
    "two_indicator_id" BIGINT NOT NULL,
    "two_indicator_name" VARCHAR(50) NOT NULL,
    "score" INT NOT NULL,
    "one_indicator_id" INT NOT NULL,
    "create_time" TIMESTAMP(0) NOT NULL,
    NOT CLUSTER PRIMARY KEY("activity_id", "two_indicator_id"),
    CHECK("activity_id" >= 0)
    ,CHECK("two_indicator_id" >= 0)) STORAGE(ON "teachers_awards", CLUSTERBTR) ;

COMMENT ON COLUMN "teachers_awards"."activity_two_indicator"."activity_id" IS '活动id';
COMMENT ON COLUMN "teachers_awards"."activity_two_indicator"."two_indicator_id" IS '二级指标id';
COMMENT ON COLUMN "teachers_awards"."activity_two_indicator"."two_indicator_name" IS '二级指标名称';
COMMENT ON COLUMN "teachers_awards"."activity_two_indicator"."score" IS '分值';
COMMENT ON COLUMN "teachers_awards"."activity_two_indicator"."one_indicator_id" IS '所属一级指标';
COMMENT ON COLUMN "teachers_awards"."activity_two_indicator"."create_time" IS '创建时间';


CREATE TABLE "teachers_awards"."announcement"
(
    "announcement_id" BIGINT IDENTITY(55, 1) NOT NULL,
    "title" VARCHAR(50) NOT NULL,
    "content" TEXT NOT NULL,
    "annex" VARCHAR(255) NOT NULL,
    "user_id" VARCHAR(64) NOT NULL,
    "user_name" VARCHAR(32) NOT NULL,
    "create_time" TIMESTAMP(0) NOT NULL,
    NOT CLUSTER PRIMARY KEY("announcement_id")) STORAGE(ON "teachers_awards", CLUSTERBTR) ;

COMMENT ON TABLE "teachers_awards"."announcement" IS '公告';
COMMENT ON COLUMN "teachers_awards"."announcement"."announcement_id" IS '主键id';
COMMENT ON COLUMN "teachers_awards"."announcement"."title" IS '标题';
COMMENT ON COLUMN "teachers_awards"."announcement"."content" IS '内容';
COMMENT ON COLUMN "teachers_awards"."announcement"."annex" IS '附件链接';
COMMENT ON COLUMN "teachers_awards"."announcement"."user_id" IS '用户id';
COMMENT ON COLUMN "teachers_awards"."announcement"."user_name" IS '用户名称';
COMMENT ON COLUMN "teachers_awards"."announcement"."create_time" IS '创建时间';


CREATE OR REPLACE  INDEX "index_user_id" ON "teachers_awards"."announcement"("user_id" ASC) STORAGE(ON "teachers_awards", CLUSTERBTR) ;

CREATE TABLE "teachers_awards"."expert_auth_indicator"
(
    "user_id" VARCHAR(32) NOT NULL,
    "two_indicator_id" INT NOT NULL,
    NOT CLUSTER PRIMARY KEY("user_id", "two_indicator_id")) STORAGE(ON "teachers_awards", CLUSTERBTR) ;

COMMENT ON TABLE "teachers_awards"."expert_auth_indicator" IS '专家授权指标';
COMMENT ON COLUMN "teachers_awards"."expert_auth_indicator"."user_id" IS '用户id';
COMMENT ON COLUMN "teachers_awards"."expert_auth_indicator"."two_indicator_id" IS '二级指标id';


CREATE TABLE "teachers_awards"."judges_verify"
(
    "user_indicator_pass_id" BIGINT IDENTITY(266, 1) NOT NULL,
    "user_activity_indicator_id" DECIMAL(20,0) NOT NULL,
    "user_activity_id" BIGINT NOT NULL,
    "judges_id" VARCHAR(32) NOT NULL,
    "judges_name" VARCHAR(32) NOT NULL,
    "judges_type" INT NOT NULL,
    "is_pass" INT NOT NULL,
    "score" INT NOT NULL,
    "opinion" VARCHAR(255) NOT NULL,
    "create_time" TIMESTAMP(0) NOT NULL,
    "delete_time" TIMESTAMP(0),
    "delete_at" DECIMAL(20,0) NOT NULL,
    NOT CLUSTER PRIMARY KEY("user_indicator_pass_id"),
    CONSTRAINT "index_unique" UNIQUE("user_activity_indicator_id", "judges_id", "judges_type", "delete_at"),
    CHECK("user_activity_indicator_id" >= 0)
    ,CHECK("judges_type" >= 0)
    ,CHECK("is_pass" >= 0)
    ,CHECK("delete_at" >= 0)) STORAGE(ON "teachers_awards", CLUSTERBTR) ;

COMMENT ON TABLE "teachers_awards"."judges_verify" IS '评委审核';
COMMENT ON COLUMN "teachers_awards"."judges_verify"."user_indicator_pass_id" IS '主键id';
COMMENT ON COLUMN "teachers_awards"."judges_verify"."user_activity_indicator_id" IS '用户活动申报二级指标id';
COMMENT ON COLUMN "teachers_awards"."judges_verify"."user_activity_id" IS '用户活动申报id';
COMMENT ON COLUMN "teachers_awards"."judges_verify"."judges_id" IS '评委id';
COMMENT ON COLUMN "teachers_awards"."judges_verify"."judges_name" IS '评委姓名';
COMMENT ON COLUMN "teachers_awards"."judges_verify"."judges_type" IS '评委类型，1：学校，2：专家，3：教育局';
COMMENT ON COLUMN "teachers_awards"."judges_verify"."is_pass" IS '0：未通过，1：通过';
COMMENT ON COLUMN "teachers_awards"."judges_verify"."score" IS '得分';
COMMENT ON COLUMN "teachers_awards"."judges_verify"."opinion" IS '审核意见';
COMMENT ON COLUMN "teachers_awards"."judges_verify"."create_time" IS '创建时间';
COMMENT ON COLUMN "teachers_awards"."judges_verify"."delete_time" IS '删除时间';
COMMENT ON COLUMN "teachers_awards"."judges_verify"."delete_at" IS '删除时间戳（毫秒）';


CREATE TABLE "teachers_awards"."notice"
(
    "notice_id" BIGINT IDENTITY(4, 1) NOT NULL,
    "notice_type" INT NOT NULL,
    "content" VARCHAR(1500) NOT NULL,
    "user_id" VARCHAR(64) NOT NULL,
    "is_read" INT NOT NULL,
    "create_time" TIMESTAMP(0) NOT NULL,
    NOT CLUSTER PRIMARY KEY("notice_id"),
    CHECK("notice_type" >= 0)
    ,CHECK("is_read" >= 0)) STORAGE(ON "teachers_awards", CLUSTERBTR) ;

COMMENT ON TABLE "teachers_awards"."notice" IS '消息通知';
COMMENT ON COLUMN "teachers_awards"."notice"."notice_id" IS '主键';
COMMENT ON COLUMN "teachers_awards"."notice"."notice_type" IS '通知类型，1：系统通知';
COMMENT ON COLUMN "teachers_awards"."notice"."content" IS '通知内容';
COMMENT ON COLUMN "teachers_awards"."notice"."user_id" IS '接收者id';
COMMENT ON COLUMN "teachers_awards"."notice"."is_read" IS '0：未读，1：已读';
COMMENT ON COLUMN "teachers_awards"."notice"."create_time" IS '创建时间';


CREATE TABLE "teachers_awards"."one_indicator"
(
    "one_indicator_id" BIGINT IDENTITY(122, 1) NOT NULL,
    "one_indicator_name" VARCHAR(50) NOT NULL,
    "content" VARCHAR(1500) NOT NULL,
    "create_time" TIMESTAMP(0) NOT NULL,
    "update_time" TIMESTAMP(0),
    "delete_time" TIMESTAMP(0),
    NOT CLUSTER PRIMARY KEY("one_indicator_id")) STORAGE(ON "teachers_awards", CLUSTERBTR) ;

COMMENT ON TABLE "teachers_awards"."one_indicator" IS '一级指标';
COMMENT ON COLUMN "teachers_awards"."one_indicator"."one_indicator_id" IS '主键id';
COMMENT ON COLUMN "teachers_awards"."one_indicator"."one_indicator_name" IS '一级指标名称';
COMMENT ON COLUMN "teachers_awards"."one_indicator"."content" IS '评分标准说明';
COMMENT ON COLUMN "teachers_awards"."one_indicator"."create_time" IS '创建时间';
COMMENT ON COLUMN "teachers_awards"."one_indicator"."update_time" IS '更新时间';
COMMENT ON COLUMN "teachers_awards"."one_indicator"."delete_time" IS '删除时间';


CREATE TABLE "teachers_awards"."operation_record"
(
    "operation_id" BIGINT IDENTITY(278, 1) NOT NULL,
    "relational_id" DECIMAL(20,0) NOT NULL,
    "relational_type" INT NOT NULL,
    "user_id" VARCHAR(64) NOT NULL,
    "user_name" VARCHAR(32) NOT NULL,
    "operation_role" INT NOT NULL,
    "operation_type" INT NOT NULL,
    "description" VARCHAR(255) NOT NULL,
    "create_time" TIMESTAMP(0) NOT NULL,
    NOT CLUSTER PRIMARY KEY("operation_id"),
    CHECK("relational_id" >= 0)
    ,CHECK("relational_type" >= 0)
    ,CHECK("operation_role" >= 0)
    ,CHECK("operation_type" >= 0)) STORAGE(ON "teachers_awards", CLUSTERBTR) ;

COMMENT ON COLUMN "teachers_awards"."operation_record"."operation_id" IS '操作id';
COMMENT ON COLUMN "teachers_awards"."operation_record"."relational_id" IS '关联的id';
COMMENT ON COLUMN "teachers_awards"."operation_record"."relational_type" IS '关联表，1：user_activity_indicator';
COMMENT ON COLUMN "teachers_awards"."operation_record"."user_id" IS '用户id';
COMMENT ON COLUMN "teachers_awards"."operation_record"."user_name" IS '用户姓名';
COMMENT ON COLUMN "teachers_awards"."operation_record"."operation_role" IS '角色，1：学校，2：专家，3：教育局，4：教师';
COMMENT ON COLUMN "teachers_awards"."operation_record"."operation_type" IS '操作类型，1：添加，2：修改，3：删除';
COMMENT ON COLUMN "teachers_awards"."operation_record"."description" IS '操作说明';
COMMENT ON COLUMN "teachers_awards"."operation_record"."create_time" IS '创建时间';


CREATE OR REPLACE  INDEX "index_relational" ON "teachers_awards"."operation_record"("relational_id" ASC,"relational_type" ASC) STORAGE(ON "teachers_awards", CLUSTERBTR) ;

CREATE TABLE "teachers_awards"."read"
(
    "id" BIGINT IDENTITY(70, 1) NOT NULL,
    "read_id" BIGINT NOT NULL,
    "read_type" INT NOT NULL,
    "read_user_id" VARCHAR(64) NOT NULL,
    "create_time" TIMESTAMP(0) NOT NULL,
    NOT CLUSTER PRIMARY KEY("id"),
    CONSTRAINT "unique_index" UNIQUE("read_id", "read_type", "read_user_id"),
    CHECK("read_id" >= 0)
    ,CHECK("read_type" >= 0)) STORAGE(ON "teachers_awards", CLUSTERBTR) ;

COMMENT ON TABLE "teachers_awards"."read" IS '记录用户的公告已读';
COMMENT ON COLUMN "teachers_awards"."read"."id" IS '主键id';
COMMENT ON COLUMN "teachers_awards"."read"."read_id" IS '已读id';
COMMENT ON COLUMN "teachers_awards"."read"."read_type" IS '1：公告';
COMMENT ON COLUMN "teachers_awards"."read"."read_user_id" IS '已读用户id';
COMMENT ON COLUMN "teachers_awards"."read"."create_time" IS '创建时间';


CREATE TABLE "teachers_awards"."two_indicator"
(
    "two_indicator_id" BIGINT IDENTITY(132, 1) NOT NULL,
    "two_indicator_name" VARCHAR(50) NOT NULL,
    "score" INT NOT NULL,
    "one_indicator_id" INT NOT NULL,
    "create_time" TIMESTAMP(0) NOT NULL,
    "update_time" TIMESTAMP(0),
    "delete_time" TIMESTAMP(0),
    NOT CLUSTER PRIMARY KEY("two_indicator_id")) STORAGE(ON "teachers_awards", CLUSTERBTR) ;

COMMENT ON TABLE "teachers_awards"."two_indicator" IS '二级指标';
COMMENT ON COLUMN "teachers_awards"."two_indicator"."two_indicator_id" IS '主键id';
COMMENT ON COLUMN "teachers_awards"."two_indicator"."two_indicator_name" IS '二级指标名称';
COMMENT ON COLUMN "teachers_awards"."two_indicator"."score" IS '分值';
COMMENT ON COLUMN "teachers_awards"."two_indicator"."one_indicator_id" IS '所属一级指标';
COMMENT ON COLUMN "teachers_awards"."two_indicator"."create_time" IS '创建时间';
COMMENT ON COLUMN "teachers_awards"."two_indicator"."update_time" IS '更新时间';
COMMENT ON COLUMN "teachers_awards"."two_indicator"."delete_time" IS '删除时间';


CREATE TABLE "teachers_awards"."user_activity"
(
    "user_activity_id" BIGINT IDENTITY(38, 1) NOT NULL,
    "activity_id" BIGINT NOT NULL,
    "user_id" VARCHAR(32) NOT NULL,
    "user_name" VARCHAR(20) NOT NULL,
    "user_sex" INT NOT NULL,
    "birthday" VARCHAR(255) NOT NULL,
    "identity_card" VARCHAR(255) NOT NULL,
    "phone" VARCHAR(255) NOT NULL,
    "subject_code" VARCHAR(20) NOT NULL,
    "school_id" VARCHAR(32) NOT NULL,
    "school_name" VARCHAR(50) NOT NULL,
    "declare_type" INT NOT NULL,
    "final_score" INT NOT NULL,
    "rank" BIGINT NOT NULL,
    "rank_prize" TINYINT NOT NULL,
    "prize" BIGINT NOT NULL,
    "create_time" TIMESTAMP(0) NOT NULL,
    "update_time" TIMESTAMP(0),
    "delete_time" TIMESTAMP(0),
    "delete_at" DECIMAL(20,0) NOT NULL,
    NOT CLUSTER PRIMARY KEY("user_activity_id"),
    CONSTRAINT "index_user_activity" UNIQUE("user_id", "activity_id", "delete_at"),
    CHECK("activity_id" >= 0)
    ,CHECK("user_sex" >= 0)
    ,CHECK("declare_type" >= 0)
    ,CHECK("final_score" >= 0)
    ,CHECK("rank" >= 0)
    ,CHECK("prize" >= 0)
    ,CHECK("delete_at" >= 0)) STORAGE(ON "teachers_awards", CLUSTERBTR) ;

COMMENT ON TABLE "teachers_awards"."user_activity" IS '用户申报活动';
COMMENT ON COLUMN "teachers_awards"."user_activity"."user_activity_id" IS '主键id';
COMMENT ON COLUMN "teachers_awards"."user_activity"."activity_id" IS '活动id';
COMMENT ON COLUMN "teachers_awards"."user_activity"."user_id" IS '用户id';
COMMENT ON COLUMN "teachers_awards"."user_activity"."user_name" IS '用户姓名';
COMMENT ON COLUMN "teachers_awards"."user_activity"."user_sex" IS '1：男，2：女';
COMMENT ON COLUMN "teachers_awards"."user_activity"."birthday" IS '出生日期';
COMMENT ON COLUMN "teachers_awards"."user_activity"."identity_card" IS '身份证号';
COMMENT ON COLUMN "teachers_awards"."user_activity"."phone" IS '手机号';
COMMENT ON COLUMN "teachers_awards"."user_activity"."subject_code" IS '科目code';
COMMENT ON COLUMN "teachers_awards"."user_activity"."school_id" IS '学校id';
COMMENT ON COLUMN "teachers_awards"."user_activity"."school_name" IS '学校名称';
COMMENT ON COLUMN "teachers_awards"."user_activity"."declare_type" IS '1:幼教、2:小学、3:初中、4:高中、5:职校、6:教研';
COMMENT ON COLUMN "teachers_awards"."user_activity"."final_score" IS '最终得分（各项通过的审核）';
COMMENT ON COLUMN "teachers_awards"."user_activity"."rank" IS '排名，如果不为0：表示为评定结果';
COMMENT ON COLUMN "teachers_awards"."user_activity"."rank_prize" IS '0：无，1：一等奖，2：二等奖，3：三等奖';
COMMENT ON COLUMN "teachers_awards"."user_activity"."prize" IS '奖金';
COMMENT ON COLUMN "teachers_awards"."user_activity"."create_time" IS '创建时间';
COMMENT ON COLUMN "teachers_awards"."user_activity"."update_time" IS '更新时间';
COMMENT ON COLUMN "teachers_awards"."user_activity"."delete_time" IS '删除时间';
COMMENT ON COLUMN "teachers_awards"."user_activity"."delete_at" IS '删除时间戳（毫秒）';


CREATE OR REPLACE  INDEX "index_activity" ON "teachers_awards"."user_activity"("activity_id" ASC) STORAGE(ON "teachers_awards", CLUSTERBTR) ;

CREATE TABLE "teachers_awards"."user_activity_indicator"
(
    "user_activity_indicator_id" BIGINT IDENTITY(100, 1) NOT NULL,
    "user_activity_id" DECIMAL(20,0) NOT NULL,
    "two_indicator_id" BIGINT NOT NULL,
    "award_date" DATE NOT NULL,
    "certificate_type" TINYINT NOT NULL,
    "certificate_url" VARCHAR(255) NOT NULL,
    "certificate_start_date" DATE,
    "certificate_end_date" DATE,
    "status" TINYINT NOT NULL,
    "finish_review_num" INT NOT NULL,
    "review_process" TINYINT NOT NULL,
    "create_time" TIMESTAMP(0) NOT NULL,
    "update_time" TIMESTAMP(0),
    "delete_time" TIMESTAMP(0),
    "delete_at" DECIMAL(20,0) NOT NULL,
    NOT CLUSTER PRIMARY KEY("user_activity_indicator_id"),
    CHECK("user_activity_id" >= 0)
    ,CHECK("two_indicator_id" >= 0)
    ,CHECK("finish_review_num" >= 0)
    ,CHECK("delete_at" >= 0)) STORAGE(ON "teachers_awards", CLUSTERBTR) ;

COMMENT ON TABLE "teachers_awards"."user_activity_indicator" IS '用户申报的活动下的指标';
COMMENT ON COLUMN "teachers_awards"."user_activity_indicator"."user_activity_indicator_id" IS '主键id';
COMMENT ON COLUMN "teachers_awards"."user_activity_indicator"."user_activity_id" IS '用户申报活动的id';
COMMENT ON COLUMN "teachers_awards"."user_activity_indicator"."two_indicator_id" IS '二级指标id';
COMMENT ON COLUMN "teachers_awards"."user_activity_indicator"."award_date" IS '获奖日期';
COMMENT ON COLUMN "teachers_awards"."user_activity_indicator"."certificate_type" IS '1.证书；需要填写证书有效期，2.证明；不需要有效期';
COMMENT ON COLUMN "teachers_awards"."user_activity_indicator"."certificate_url" IS '证书url';
COMMENT ON COLUMN "teachers_awards"."user_activity_indicator"."certificate_start_date" IS '证书有效期——开始时间';
COMMENT ON COLUMN "teachers_awards"."user_activity_indicator"."certificate_end_date" IS '证书有效期——结束时间';
COMMENT ON COLUMN "teachers_awards"."user_activity_indicator"."status" IS '审核状态说明：
1.已提交：教师提交申报、进入审核流程，均显示已提交状态
2.已通过：仅当教育局通过最终审核，显示已通过，无法进行修改、撤销操作
3.学校未通过：学校审核不通过，审核流程不进入专家审核阶段
4.专家未通过：无论专家审核是否通过，均进入教育局审核阶段
5.教育局未通过：该项不通过
6.结果已评定，未审批';
COMMENT ON COLUMN "teachers_awards"."user_activity_indicator"."finish_review_num" IS '当前已审核人数';
COMMENT ON COLUMN "teachers_awards"."user_activity_indicator"."review_process" IS '当前审核进程，1：学校，2：专家，3：教育局，4：结束';
COMMENT ON COLUMN "teachers_awards"."user_activity_indicator"."create_time" IS '创建时间';
COMMENT ON COLUMN "teachers_awards"."user_activity_indicator"."update_time" IS '更新时间';
COMMENT ON COLUMN "teachers_awards"."user_activity_indicator"."delete_time" IS '删除时间';
COMMENT ON COLUMN "teachers_awards"."user_activity_indicator"."delete_at" IS '删除时间戳（毫秒）';


CREATE OR REPLACE  INDEX "index_user_activity_id" ON "teachers_awards"."user_activity_indicator"("user_activity_id" ASC) STORAGE(ON "teachers_awards", CLUSTERBTR) ;

CREATE TABLE "teachers_awards"."user_info"
(
    "user_id" VARCHAR(32) NOT NULL,
    "user_name" VARCHAR(32) NOT NULL,
    "user_sex" INT NOT NULL,
    "birthday" VARCHAR(255) NOT NULL,
    "identity_card" VARCHAR(255) NOT NULL,
    "phone" VARCHAR(255) NOT NULL,
    "avatar" VARCHAR(255) NOT NULL,
    "subject_code" VARCHAR(20) NOT NULL,
    "school_id" VARCHAR(32) NOT NULL,
    "school_name" VARCHAR(50) NOT NULL,
    "declare_type" INT NOT NULL,
    "role" INT NOT NULL,
    "export_auth" TINYINT NOT NULL,
    "auth_day" DATE,
    "year" INT NOT NULL,
    "create_time" TIMESTAMP(0) NOT NULL,
    "update_time" TIMESTAMP(0),
    NOT CLUSTER PRIMARY KEY("user_id"),
    CHECK("user_sex" >= 0)
    ,CHECK("declare_type" >= 0)) STORAGE(ON "teachers_awards", CLUSTERBTR) ;

COMMENT ON TABLE "teachers_awards"."user_info" IS '记录用户当前基本信息';
COMMENT ON COLUMN "teachers_awards"."user_info"."user_id" IS '用户id';
COMMENT ON COLUMN "teachers_awards"."user_info"."user_name" IS '用户名称';
COMMENT ON COLUMN "teachers_awards"."user_info"."user_sex" IS '1：男，2：女';
COMMENT ON COLUMN "teachers_awards"."user_info"."birthday" IS '出生日期';
COMMENT ON COLUMN "teachers_awards"."user_info"."identity_card" IS '身份证号';
COMMENT ON COLUMN "teachers_awards"."user_info"."phone" IS '手机号';
COMMENT ON COLUMN "teachers_awards"."user_info"."avatar" IS '头像';
COMMENT ON COLUMN "teachers_awards"."user_info"."subject_code" IS '科目code';
COMMENT ON COLUMN "teachers_awards"."user_info"."school_id" IS '学校id';
COMMENT ON COLUMN "teachers_awards"."user_info"."school_name" IS '学校名称';
COMMENT ON COLUMN "teachers_awards"."user_info"."declare_type" IS '1:幼教、2:小学、3:初中、4:高中、5:职校、6:教研';
COMMENT ON COLUMN "teachers_awards"."user_info"."role" IS '角色，1：学校，2：专家，4：教育局，8：教师，16：超级管理员，多个则相加';
COMMENT ON COLUMN "teachers_awards"."user_info"."export_auth" IS '专家是否授权 1：未授权 2：已授权';
COMMENT ON COLUMN "teachers_awards"."user_info"."auth_day" IS '授权日期';
COMMENT ON COLUMN "teachers_awards"."user_info"."year" IS '年份';
COMMENT ON COLUMN "teachers_awards"."user_info"."create_time" IS '创建时间';
COMMENT ON COLUMN "teachers_awards"."user_info"."update_time" IS '更新时间';


