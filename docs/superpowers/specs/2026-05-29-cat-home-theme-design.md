# 首页黑白猫猫主题设计规格

## 背景

当前默认首页位于 `web/default/src/features/home`，由首页路由 `web/default/src/routes/index.tsx` 渲染。首页在没有后台自定义内容时展示默认落地页；如果后台配置了自定义 Markdown 或 URL，则直接展示自定义内容。本次改造只针对默认首页，不改变自定义首页、认证、导航、文档链接、定价页、控制台入口和 API Demo 轮播的数据逻辑。

用户选择方案 2：重排首页为更强的整屏插画叙事，主题是“可爱的猫猫在玩耍”，黑白配色，同时不影响现有功能。

## 目标

- 将默认首页从蓝紫科技渐变风改为黑白漫画纸张风。
- 第一屏建立明确的“猫猫在 API 游乐场玩耍”视觉记忆点。
- 保留所有现有功能入口和条件渲染逻辑。
- 保持响应式、可访问性和暗色模式可读性。
- 尽量把改动限制在 `features/home` 及首页专用静态资产范围。

## 非目标

- 不修改后台、认证、计费、渠道、模型、控制台等业务逻辑。
- 不移除或替换项目受保护的品牌、作者、版权和许可证信息。
- 不改变公共导航、页脚链接、路由结构或登录后跳转行为。
- 不新增全局主题预设，不影响控制台页面的视觉系统。
- 不修改后台自定义首页内容的 Markdown/iframe 分支。

## 视觉方向

首页采用黑、白、灰三色为主，少量纸张纹理、墨线、网点和手绘边框。主视觉是一只黑白猫猫在由“API 线缆、模型标签、终端窗口、毛线团”组成的游乐场中玩耍。猫猫负责传递可爱和亲近感，终端/API 元素负责保留产品身份。

第一屏不再是左右常规文本加卡片布局，而是更叙事化的舞台：

- 背景像一张干净的漫画纸，带轻微灰阶纹理和细网格。
- 主标题、按钮和 API Demo 仍然可操作，但围绕猫猫游乐场排列。
- 猫猫、毛线/API 线缆、终端卡片之间用黑色手绘线连接，形成“请求被猫猫玩着路由出去”的趣味隐喻。
- 统计区、功能区、步骤区和 CTA 延续黑白纸片、猫爪、线团、贴纸标签的视觉语言。

## 资产策略

主视觉使用项目内保存的生成式位图资产，风格要求为黑白、可爱、干净、适合网页首屏。资产保存到 `web/default/src/assets/home/`，由首页组件引用，避免依赖外部网络资源。配套的小装饰元素可用 CSS 或轻量 SVG 实现，例如猫爪、线团、贴纸边框和手绘连线。

如果生成图不够稳定，回退方案是在首页组件内用 SVG/CSS 构建简化黑白猫猫插画，以保证构建和部署可控。

## 信息架构

默认首页保留现有五段内容，但重新组织视觉层级：

1. **Hero：猫猫 API 游乐场**
   - 保留登录态按钮分支：已登录显示进入控制台和文档；未登录显示开始使用、查看价格、文档。
   - 保留 `useStatus()` 读取到的文档链接逻辑。
   - 保留 `HeroTerminalDemo`，但将外观改成黑白贴纸/漫画终端风格。
   - 标题文案继续使用现有 i18n key，不新增硬编码中文。

2. **Stats：猫爪计数条**
   - 保留计数动画和四个统计项。
   - 视觉改成黑白分栏纸条或猫爪印计数。

3. **Features：黑白贴纸功能板**
   - 保留现有核心功能文案和图标语义。
   - Bento 区域变成贴纸、便签、漫画框组合。
   - 避免过多彩色图标，图标统一为黑白线稿。

4. **How It Works：三步猫猫路线图**
   - 保留配置、连接、监控三步。
   - 用猫猫脚印或毛线轨迹串联步骤。

5. **CTA：猫猫把线团推向入口**
   - 未登录用户仍展示注册和价格按钮。
   - 已登录用户继续不展示 CTA，保持当前行为。

## 组件边界

优先修改这些文件：

- `web/default/src/features/home/components/sections/hero.tsx`
- `web/default/src/features/home/components/hero-terminal-demo.tsx`
- `web/default/src/features/home/components/sections/stats.tsx`
- `web/default/src/features/home/components/sections/features.tsx`
- `web/default/src/features/home/components/sections/how-it-works.tsx`
- `web/default/src/features/home/components/sections/cta.tsx`

可新增首页专用组件：

- `web/default/src/features/home/components/cat-playground-visual.tsx`
- `web/default/src/features/home/components/cat-doodle-elements.tsx`

可新增首页专用资产目录：

- `web/default/src/assets/home/`

不修改 `Home` 顶层内容分支，除非只是为默认首页包裹一个主题 class，且不影响自定义内容分支。

## 数据和状态流

数据流保持现状：

- `Home` 读取 `useAuthStore()` 判断是否已登录。
- `Home` 读取 `useHomePageContent()`；有自定义内容时直接返回自定义内容。
- `Hero` 读取 `useStatus()` 获取 `docs_link`。
- `HeroTerminalDemo` 自己维护 tab 轮播和点击切换状态。
- 各区块继续使用 `useTranslation()` 和现有 i18n key。

本次不新增 API 请求、全局 store、路由参数或持久化状态。

## 交互和动效

- Hero 首屏可以有轻微入场动效：标题、按钮、猫猫主视觉和终端贴纸错峰出现。
- 猫猫或线团可以做非常轻微的 CSS 动效，例如尾巴摆动、线团滚动暗示；必须遵守 `prefers-reduced-motion`。
- API Demo 轮播继续保留现有点击和自动轮播行为。
- 按钮 hover、focus、active 状态沿用项目 Button 组件能力，并适配黑白视觉。

## 响应式

- 桌面端：第一屏可以使用重排舞台式布局，猫猫主视觉占据视觉中心，标题和终端贴纸围绕布局。
- 平板端：主视觉与文本上下错落，终端卡片保持可读。
- 移动端：标题、按钮、主视觉、API Demo 纵向排列；主视觉裁切不得遮挡按钮和正文。
- 所有按钮文字和 API Demo 文本不得溢出容器。

## 可访问性

- 主视觉如果只是装饰，使用空 `alt` 或 `aria-hidden`。
- 保留标题层级和语义化 section。
- 黑白配色下正文和按钮对比度满足 WCAG AA。
- 可点击元素保留键盘焦点态。
- 动效尊重系统减少动态效果设置。

## 测试与验证

实现后至少执行：

- `bun run typecheck`
- `bun run build`
- 使用本地开发服务器检查首页。
- 用浏览器验证桌面和移动视口：首屏不空白、不遮挡、按钮可点、API Demo 可切换。
- 验证后台自定义首页内容存在时，默认猫猫首页不会渲染，仍展示自定义 Markdown 或 iframe。

## 风险和缓解

- **生成图风格不稳定**：先生成并筛选一张干净黑白主视觉；若不理想，回退到 SVG/CSS 插画。
- **重排 Hero 影响入口可见性**：按钮区优先级高于装饰层，移动端先保证操作路径。
- **黑白风过于单调**：使用纸张纹理、线稿密度、阴影层次和贴纸形状增加区分度，而不是引入彩色渐变。
- **全局视觉被误伤**：不改全局主题变量，不改控制台布局，只在首页组件内使用局部 class。
