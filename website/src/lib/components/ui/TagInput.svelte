<script lang="ts">
    let { tags = $bindable([]), max = 10, placeholder = "Enter tags...", title = "Your Tags" } = $props();

    let currentInput = $state('');
    let errorText = $state('');


    function addTag(e: KeyboardEvent) {
		if (e.key === 'Enter' && currentInput.trim()) {
			e.preventDefault();
			if (tags.length >= max) {
				errorText = 'Maximum of ten tags';
				return;
			}
			// 如果标签不存在则添加
			if (!tags.includes(currentInput.trim())) {
				tags.push(currentInput.trim());
			}
			currentInput = ''; // 清空输入
		}
	}
</script>

<div class="flex w-full max-w-sm flex-col">
	<div class="label"><span class="label-text font-medium">{title}</span></div>

	<input
		type="text"
		placeholder="{placeholder}"
		class="input-bordered input w-full"
		bind:value={currentInput}
		onkeydown={addTag}
	/>

	<div class="mt-2 flex flex-wrap gap-2">
		{#each tags as tag, i}
			<div class="badge gap-2 badge-soft p-3 badge-primary">
				{tag}
				<button type="button" onclick={() => tags.splice(i, 1)} class="text-xs">✕</button>
			</div>
		{/each}
	</div>
    <p class="text-red-500">{errorText}</p>
</div>
