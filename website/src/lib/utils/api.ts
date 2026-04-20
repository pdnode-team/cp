/**
 * 转换器：将普通对象转换为 PocketBase 兼容的 FormData
 */
export function toFormData(data: Record<string, any>): FormData {
    const formData = new FormData();

    for (const [key, value] of Object.entries(data)) {
        // 跳过空值，除非你想显式清空字段
        if (value === undefined || value === null) continue;

        if (value instanceof FileList) {
            // 处理多文件
            Array.from(value).forEach(file => formData.append(key, file));
        } else if (value instanceof File) {
            // 处理单文件
            formData.append(key, value);
        } else if (Array.isArray(value)) {
            // 处理数组 (tags, 关系 ID 等)
            // 如果数组为空，append 一个空字符串告诉 PB 清空该字段
            if (value.length === 0) {
                formData.append(key, '');
            } else {
                value.forEach(item => formData.append(key, item));
            }
        } else if (typeof value === 'boolean') {
            // PocketBase 接受 "true" / "false" 字符串
            formData.append(key, value.toString());
        } else {
            // 其他基本类型 (string, number)
            formData.append(key, value.toString());
        }
    }

    return formData;
}