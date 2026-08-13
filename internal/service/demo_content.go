package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/gorm"
)

type demoCategory struct {
	name        string
	slug        string
	description string
}

type demoPost struct {
	title       string
	slug        string
	description string
	category    string
	content     string
}

var demoCategories = []demoCategory{
	{name: "前端工程", slug: "frontend", description: "浏览器、交互与现代前端工程实践"},
	{name: "后端开发", slug: "backend", description: "服务端架构、数据与接口设计"},
	{name: "工程随笔", slug: "engineering", description: "开发过程中的思考与复盘"},
	{name: "生活记录", slug: "life", description: "工作之外值得留下的片段"},
}

var demoPosts = []demoPost{
	{title: "从一次慢查询开始理解数据库索引", slug: "demo-database-index", description: "用一个真实查询过程梳理联合索引、回表和执行计划。", category: "backend", content: "<h2>问题从哪里开始</h2><p>当数据量逐渐增长，同一条查询可能从几毫秒变成几秒。排查的第一步不是立刻增加索引，而是先确认过滤、排序和返回字段。</p><blockquote>索引不是越多越好，它应该服务于稳定而高频的访问路径。</blockquote><h2>观察执行计划</h2><p>通过执行计划可以判断是否命中索引、扫描了多少行，以及排序是否产生了额外开销。</p>"},
	{title: "React 状态设计中的三个边界", slug: "demo-react-state-boundaries", description: "组件状态、上下文和服务端状态分别应该解决什么问题。", category: "frontend", content: "<h2>状态离使用位置越近越好</h2><p>只影响一个表单的值没有必要进入全局状态。组件状态负责短生命周期交互，上下文负责跨层共享，而服务端状态应该由请求缓存负责。</p><h2>避免复制服务端数据</h2><p>当同一份数据被复制到多个状态容器中，刷新和失效会变得难以推断。</p>"},
	{title: "为博客设计一条可靠的发布链路", slug: "demo-publishing-pipeline", description: "从本地草稿、服务端草稿到正式发布的状态划分。", category: "engineering", content: "<h2>草稿并不只有一种</h2><p>浏览器快照用于防止意外丢失，服务端草稿用于跨设备继续编辑，正式内容则面向读者。三者职责不同，不应该共用一个模糊状态。</p><h2>发布是一次明确的状态转换</h2><p>发布动作需要校验正文、可见范围和关联数据，并在成功后清理本地快照。</p>"},
	{title: "Hertz 中间件如何贯穿一次请求", slug: "demo-hertz-middleware", description: "用认证链路理解中间件、上下文和控制器之间的分工。", category: "backend", content: "<h2>中间件处理横切逻辑</h2><p>认证、日志和限流都不是单个业务接口独有的能力。中间件先解析凭证，再把用户身份写入请求上下文。</p><pre><code>UseAuth(false) -&gt; Controller -&gt; Service</code></pre><p>服务层仍然需要执行资源级权限判断，因为登录并不等于拥有目标内容的访问权。</p>"},
	{title: "让页面动起来，但不要打断阅读", slug: "demo-motion-and-reading", description: "动效应该帮助理解层级，而不是抢走内容注意力。", category: "frontend", content: "<h2>动效需要有原因</h2><p>进入、展开和状态反馈是三类常见用途。持续运动的装饰应该保持克制，并尊重系统的减少动态效果设置。</p><h2>稳定布局优先</h2><p>动画元素应有稳定尺寸，避免加载后推动标题、按钮或正文发生位移。</p>"},
	{title: "周末散步时记下的几件小事", slug: "demo-weekend-notes", description: "离开屏幕之后，重新观察城市里缓慢发生的变化。", category: "life", content: "<h2>慢一点走</h2><p>熟悉的街道在不赶时间时会露出很多细节：新开的旧书店、窗台上的植物，还有傍晚逐渐变长的影子。</p><p>记录并不一定需要宏大的主题，能准确保留当时的感受就已经足够。</p>"},
}

func EnsureDemoContent(ctx context.Context) error {
	rootUser, err := repo.GetRootUser(ctx)
	if err != nil {
		return fmt.Errorf("get demo content author: %w", err)
	}

	categories := make(map[string]model.Category, len(demoCategories))
	for _, item := range demoCategories {
		category, err := repo.GetCategoryBySlug(ctx, item.slug)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			category = model.Category{Name: item.name, Slug: item.slug, Description: item.description}
			if err := repo.CreateCategory(ctx, &category); err != nil {
				return fmt.Errorf("create demo category %s: %w", item.slug, err)
			}
		} else if err != nil {
			return fmt.Errorf("get demo category %s: %w", item.slug, err)
		}
		categories[item.slug] = category
	}

	for index, item := range demoPosts {
		if _, err := repo.GetPostBySlug(ctx, item.slug); err == nil {
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("get demo post %s: %w", item.slug, err)
		}

		category := categories[item.category]
		publishedAt := time.Now().Add(-time.Duration(index+1) * 24 * time.Hour)
		post := model.Post{
			PostBase: model.PostBase{
				Title:        item.title,
				Content:      item.content,
				DraftContent: item.content,
				Description:  item.description,
				Type:         "article",
				Slug:         item.slug,
				CategoryID:   &category.ID,
				Visibility:   constant.PostVisibilityPublic,
				PublishedAt:  &publishedAt,
			},
			AuthorID: rootUser.ID,
		}
		if err := repo.CreatePost(ctx, &post); err != nil {
			return fmt.Errorf("create demo post %s: %w", item.slug, err)
		}
	}

	return nil
}
