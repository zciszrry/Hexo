hexo.extend.filter.register('before_post_render', function(data) {
  // 更健壮的正则，匹配各种图片语法
  const markdownImageRegex = /!\[\s*(.*?)\s*\]\(\s*(.*?)\s*\)/g;
  let changedCount = 0;
  
  data.content = data.content.replace(markdownImageRegex, function(match, altText, fullPath) {
    // 跳过已经是web URL或数据URI的图片
    if (fullPath.startsWith('http') || fullPath.startsWith('data:') || fullPath.startsWith('/')) {
      return match;
    }
    
    // 处理空 alt 文本
    const cleanAltText = altText || '';
    
    // 更健壮的文件名提取（处理Windows和Unix路径）
    const normalizedPath = fullPath.replace(/\\/g, '/');
    const pathParts = normalizedPath.split('/');
    let fileName = pathParts[pathParts.length - 1];
    
    // 移除可能的URL参数或锚点
    fileName = fileName.split('?')[0].split('#')[0];
    
    // 检查是否需要替换（路径包含目录且不是当前目录引用）
    const needsReplacement = pathParts.length > 1 && 
                            !(pathParts.length === 2 && pathParts[0] === '.');
    
    if (needsReplacement) {
      changedCount++;
      console.log(`[Hexo Image Filter] 替换图片路径: "${fullPath}" -> "${fileName}"`);
      return `![${cleanAltText}](${fileName})`;
    }
    
    return match;
  });
  
  if (changedCount > 0) {
    console.log(`[Hexo Image Filter] 在文章 "${data.title}" 中替换了 ${changedCount} 个图片路径`);
  }
  
  return data;
}, 1);