---
title: Kaggle部署Stable-Diffusion
date: 2025-09-29 10:26:50
tags: Stable-Diffusion
---

​	炼图初步。九品炼图师中期巅峰半步大圆满境。

<!--more-->

​	比如我可以炼几张粉色狐狸：

​	![](D:\Hexo\source\_posts\Kaggle部署Stable-Diffusion\14.png)



### 前言

​	AI盛行那段时间，其实一直想自己炼图来着，但也一直止步于想法。

​	这两天很闲，小研究了一下。

​	首先看了B站up主：[秋葉aaaki](https://space.bilibili.com/12566101) 的一些教程，下载了整合包，但是本人小小的笔记本感觉跑不起来的样子。

​	然后想着怎么云炼图。

​	其实按理来说并不需要自己来搭建部署，网上有云应用，比如[端脑云](https://cephalon.cloud/aigc)。而且也省事省时间。

​	不过我刚好比较闲，而且云应用好像会一些审核，所有还是想着怎么搭一个。

​	可是，真买一个云服务器而且还是需要一定配置的云服务器，太贵了呀...又没有别的需求，而且炼图只是一时兴起。

​	所以去找找免费的，详见：[整理一些免费的AI绘画云端部署](https://zhuanlan.zhihu.com/p/685049850)。

### Kaggle

​	新建NoteBoook，复制下面代码并运行。

~~~~python
# 📁 全新 Stable Diffusion WebUI 部署脚本
# 专为 Kaggle 环境优化，包含完整的模型和插件配置

import os
import subprocess
import time
import urllib.request
from pathlib import Path

# ==================== 配置区域 ====================
# 在这里修改您的配置

# 模型和插件配置
isControlNet = True  # 开启ControlNet，下载基本模型需要2-3分钟
isSadTalker = False  # 虚拟数字人插件

# 插件列表
extensions = [
    'https://github.com/Elldreth/loopback_scaler',
    'https://github.com/jexom/sd-webui-depth-lib',
    'https://github.com/AlUlkesh/stable-diffusion-webui-images-browser',
    'https://github.com/nonnonstop/sd-webui-3d-open-pose-editor',
    'https://github.com/2575044704/stable-diffusion-webui-localization-zh_CN2.git',
    'https://github.com/opparco/stable-diffusion-webui-two-shot',
    'https://github.com/DominikDoom/a1111-sd-webui-tagcomplete',
    'https://github.com/pkuliyi2015/multidiffusion-upscaler-for-automatic1111',
    'https://github.com/hnmr293/sd-webui-cutoff',
    'https://github.com/hako-mikan/sd-webui-lora-block-weight',
    'https://github.com/catppuccin/stable-diffusion-webui',
    'https://github.com/deforum-art/sd-webui-deforum',
    'https://github.com/KaggleSD/sd-extension-system-info',
    'https://github.com/continue-revolution/sd-webui-animatediff',
    'https://github.com/adieyal/sd-dynamic-prompts.git',
    'https://github.com/hako-mikan/sd-webui-supermerger',
    'https://github.com/Bing-su/adetailer',
    'https://github.com/thisjam/sd-webui-oldsix-prompt',
    'https://github.com/Echoflare/a1111-sd-encrypt-image'
]

# 模型下载链接
sd_model_urls = [
    '[C站热门|真人]麦橘v6.safetensors:https://civitai.com/api/download/models/94640',
    'Real:https://civitai.com/api/download/models/475568',
    'https://civitai.com/api/download/models/272376',
    '[萌二次元]131-half.safetensors:https://huggingface.co/datasets/ACCC1380/private-model/resolve/main/kaggle/input/museum/131-half.safetensors'
]

vae_model_urls = [
    'https://huggingface.co/datasets/VASVASVAS/vae/resolve/main/pastel-waifu-diffusion.vae.pt'
]

lora_model_urls = [
    'https://huggingface.co/amaru96vn/Add_Detail_Lora/resolve/main/add_detail.safetensors',
    'https://huggingface.co/latent-consistency/lcm-lora-sdv1-5/resolve/main/pytorch_lora_weights.safetensors'
]

controlnet_urls = [
    'https://huggingface.co/comfyanonymous/ControlNet-v1-1_fp16_safetensors/resolve/main/control_v11p_sd15_canny_fp16.safetensors',
    'https://huggingface.co/comfyanonymous/ControlNet-v1-1_fp16_safetensors/resolve/main/control_v11p_sd15_openpose_fp16.safetensors',
    'https://huggingface.co/comfyanonymous/ControlNet-v1-1_fp16_safetensors/resolve/main/control_v11p_sd15s2_lineart_anime_fp16.safetensors',
    'https://huggingface.co/comfyanonymous/ControlNet-v1-1_fp16_safetensors/resolve/main/control_v11u_sd15_tile_fp16.safetensors',
    'https://huggingface.co/DionTimmer/controlnet_qrcode-control_v1p_sd15/resolve/main/control_v1p_sd15_qrcode.safetensors'
]

embeddings_urls = [
    'https://huggingface.co/datasets/sukaka/sd_configs/resolve/main/%E4%BA%BA%E4%BD%93%E4%BF%AE%E6%AD%A3/EasyNegative.pt',
    'https://huggingface.co/datasets/sukaka/sd_configs/resolve/main/%E4%BA%BA%E4%BD%93%E4%BF%AE%E6%AD%A3/bad-artist-anime.pt',
    'https://huggingface.co/datasets/sukaka/sd_configs/resolve/main/%E4%BA%BA%E4%BD%93%E4%BF%AE%E6%AD%A3/bad-hands-5.pt'
]

# ==================== 核心函数 ====================

def run_command(cmd):
    """执行命令并返回结果"""
    try:
        result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
        if result.returncode == 0:
            return result.stdout
        else:
            print(f"❌ 命令失败: {cmd}")
            if result.stderr:
                print(f"错误: {result.stderr}")
            return None
    except Exception as e:
        print(f"❌ 执行异常: {e}")
        return None

def setup_environment():
    """设置基础环境"""
    print("🔧 设置环境...")
    
    # 更新系统包
    run_command("apt-get update -qq")
    run_command("apt-get install -qq -y wget git aria2")
    
    # 安装 Python 依赖
    run_command("pip install -q torch==2.1.1 torchvision==0.16.1 --index-url https://download.pytorch.org/whl/cu118")
    
    print("✅ 环境设置完成")

def clone_webui():
    """克隆 Stable Diffusion WebUI"""
    print("📥 克隆 Stable Diffusion WebUI...")
    
    %cd /kaggle/working
    
    # 删除旧目录（如果存在）
    if os.path.exists("stable-diffusion-webui"):
        run_command("rm -rf stable-diffusion-webui")
    
    # 克隆最新代码
    run_command("git clone https://github.com/AUTOMATIC1111/stable-diffusion-webui.git")
    
    if os.path.exists("stable-diffusion-webui"):
        print("✅ WebUI 克隆成功")
        return True
    else:
        print("❌ WebUI 克隆失败")
        return False

def install_dependencies():
    """安装项目依赖"""
    print("📦 安装依赖...")
    
    %cd /kaggle/working/stable-diffusion-webui
    
    # 安装 requirements.txt
    if os.path.exists("requirements.txt"):
        print("安装 requirements.txt...")
        run_command("pip install -q -r requirements.txt")
    
    # 安装额外依赖
    extra_deps = [
        "pytorch_lightning",
        "omegaconf", 
        "einops",
        "accelerate",
        "transformers",
        "open_clip_torch"
    ]
    
    for dep in extra_deps:
        run_command(f"pip install -q {dep}")
    
    print("✅ 依赖安装完成")

def download_file(url, dest_path):
    """下载单个文件"""
    try:
        os.makedirs(os.path.dirname(dest_path), exist_ok=True)
        
        if url.startswith('http'):
            print(f"📥 下载: {os.path.basename(dest_path)}")
            
            # 使用 aria2 多线程下载
            aria2_cmd = f'aria2c -x 16 -s 16 -k 1M "{url}" -d "{os.path.dirname(dest_path)}" -o "{os.path.basename(dest_path)}"'
            result = run_command(aria2_cmd)
            
            if result and os.path.exists(dest_path):
                print(f"✅ 下载成功: {os.path.basename(dest_path)}")
                return True
        else:
            # 处理本地文件
            if os.path.exists(url):
                run_command(f"cp '{url}' '{dest_path}'")
                return True
                
    except Exception as e:
        print(f"❌ 下载失败 {url}: {e}")
    
    return False

def download_models():
    """下载所有模型文件"""
    print("🎯 开始下载模型...")
    
    base_path = "/kaggle/working/stable-diffusion-webui"
    
    # 创建目录结构
    directories = {
        'checkpoints': f"{base_path}/models/Stable-diffusion",
        'vae': f"{base_path}/models/VAE", 
        'lora': f"{base_path}/models/Lora",
        'controlnet': f"{base_path}/extensions/sd-webui-controlnet/models",
        'embeddings': f"{base_path}/embeddings"
    }
    
    for dir_name, dir_path in directories.items():
        os.makedirs(dir_path, exist_ok=True)
        print(f"📁 创建目录: {dir_path}")
    
    # 下载模型文件
    model_groups = [
        ('SD模型', sd_model_urls, directories['checkpoints']),
        ('VAE模型', vae_model_urls, directories['vae']),
        ('Lora模型', lora_model_urls, directories['lora']),
        ('ControlNet模型', controlnet_urls, directories['controlnet']),
        ('Embeddings', embeddings_urls, directories['embeddings'])
    ]
    
    for group_name, urls, dest_dir in model_groups:
        if urls:
            print(f"\n📦 下载{group_name}...")
            for item in urls:
                if ':' in item and not item.startswith('http'):
                    # 处理带自定义文件名的链接
                    filename, url = item.split(':', 1)
                    dest_path = os.path.join(dest_dir, filename)
                else:
                    # 直接从URL提取文件名
                    url = item
                    filename = os.path.basename(url)
                    dest_path = os.path.join(dest_dir, filename)
                
                download_file(url, dest_path)
    
    print("✅ 模型下载完成")

def install_extensions():
    """安装扩展插件"""
    print("🔌 安装扩展插件...")
    
    %cd /kaggle/working/stable-diffusion-webui/extensions
    
    for extension_url in extensions:
        try:
            extension_name = os.path.basename(extension_url).replace('.git', '')
            print(f"📥 安装: {extension_name}")
            
            # 克隆扩展
            run_command(f"git clone {extension_url}")
            
        except Exception as e:
            print(f"❌ 扩展安装失败 {extension_url}: {e}")
    
    print("✅ 扩展安装完成")

def start_webui():
    """启动 WebUI"""
    # print("🚀 启动 Stable Diffusion WebUI...")
    
    # %cd /kaggle/working/stable-diffusion-webui
    
    # # 启动参数
    # launch_args = [
    #     "--share",
    #     "--enable-insecure-extension-access", 
    #     "--listen",
    #     "--skip-torch-cuda-test",
    #     "--no-half-vae",
    #     "--opt-split-attention"
    # ]
    
    # launch_cmd = f"python launch.py {' '.join(launch_args)}"
    # print(f"🎯 启动命令: {launch_cmd}")
    
    print("⏳ 启动中，请耐心等待...")
    # print("📢 成功后会显示 Gradio 公共链接")

    # print("请运行直接启动代码块")
    print("📢 成功后会显示 Gradio 公共链接")
    
    # 启动 WebUI
    #run_command(launch_cmd)
    # 启动
    
    %cd /kaggle/working
    %cd stable-diffusion-webui
    !python launch.py --share --listen --enable-insecure-extension-access --skip-torch-cuda-test

# ==================== 主执行流程 ====================

def main():
    """主执行函数"""
    print("🎯 Stable Diffusion WebUI 全新部署")
    print("=" * 60)
    
    # 检查 GPU
    print("🔍 检查 GPU...")
    run_command("nvidia-smi")
    
    # 执行部署步骤
    steps = [
        ("环境设置", setup_environment),
        ("克隆WebUI", clone_webui),
        ("安装依赖", install_dependencies),
        ("下载模型", download_models),
        ("安装扩展", install_extensions),
        ("启动WebUI", start_webui)
    ]
    
    for step_name, step_func in steps:
        print(f"\n{'='*40}")
        print(f"🔄 执行: {step_name}")
        print(f"{'='*40}")
        
        try:
            success = step_func()
            if success is False:  # 如果函数明确返回False
                print(f"❌ {step_name} 失败")
                break
            time.sleep(2)
        except Exception as e:
            print(f"❌ {step_name} 出错: {e}")
            break
    
    print("\n🎉 部署流程完成！")

# 执行主函数
if __name__ == "__main__":
    main()
~~~~

​	结果如下：

​	![](D:\Hexo\source\_posts\Kaggle部署Stable-Diffusion\11.png)

---

​	显示链接如下：

![](D:\Hexo\source\_posts\Kaggle部署Stable-Diffusion\12.png)

***

访问链接即WebUI：

![](D:\Hexo\source\_posts\Kaggle部署Stable-Diffusion\13.png)

​	具体炼图操作自行查询教程。如，[Stable Diffusion文档](https://docs.stablediffusion.cn/article/1.html)。

### Kaggle环境持久化

​	如果上面的操作能成功的话，现在会面临一个也许称得上问题的问题：Kaggle每次会清理会话环境，也就是说下次你再来还需要重新克隆仓库，下载模型，安装依赖，最后启动。

​	目前网上的办法是利用Kaggle的DataSet进行一个文件的备份和挂载，不过我没弄。。。懒，启动速度也能接受。

### 模型训练

#### 	为什么需要自己训练模型：

##### 1. 个性化需求
- **特定角色**：你的原创动漫角色、公司吉祥物
- **独特风格**：个人画风、品牌视觉风格
- **专业领域**：医疗影像、工业设计、特定产品
- **私有数据**：公司内部设计、个人照片风格

##### 2. 预训练模型的局限性

| 限制类型     | 说明                         | 例子                   |
| ------------ | ---------------------------- | ---------------------- |
| **风格限制** | 只能生成训练数据中的风格     | 无法生成你的个人画风   |
| **内容限制** | 无法生成训练数据中没有的概念 | 你的公司logo、特定产品 |
| **质量限制** | 可能无法满足专业需求         | 商业级图像质量要求     |
| **版权问题** | 某些商业用途需要自有模型     | 避免版权纠纷           |

### 总结

​	虽然很想自己训练特点的模型生成想要的图，不过一怕麻烦二以后大概率没时间（工作外的时间要留给游戏的），目前就这样叭。
