<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import pb from '$lib/pocketbase';
	import { onMount } from 'svelte';

	let charTags = $state<string[]>([]);
	let charTagsCurrentInput = $state('');
	let charErrorText = $state('');

	let record = $state<any>();

	function newCharAddTag(e: KeyboardEvent) {
		if (e.key === 'Enter' && charTagsCurrentInput.trim()) {
			e.preventDefault();
			if (charTags.length >= 10) {
				charErrorText = 'Maximum of ten tags';
				return;
			}
			// 如果标签不存在则添加
			if (!charTags.includes(charTagsCurrentInput.trim())) {
				charTags.push(charTagsCurrentInput.trim());
			}
			charTagsCurrentInput = ''; // 清空输入
		}
	}

	let charName = $state('');
	let charDescription = $state('');
	let charOrigin = $state('');
	let charPictures: FileList | undefined = $state();

	const updateCharacter = async () => {
		const formData = new FormData();
		charErrorText = '';

		if (charPictures && charPictures.length !== 0) {
			formData.delete('images');
			for (let file of charPictures) {
				formData.append('images', file);
			}
		}

		formData.append('name', charName);
		formData.append('description', charDescription);
		formData.append('origin', charOrigin);
		formData.delete('tag_names');
		charTags.forEach((tag) => {
			formData.append('tag_names', tag);
		});

		formData.append('owner', pb.authStore.record!.id);

		try {
			await pb.collection('characters').update(page.params.id!, formData);
			goto(`/characters/${page.params.id}`);
		} catch (err: any) {
			charErrorText = err.data.data?.message ?? 'Update failed. Please try again.';

			const firstKey = Object.keys(err.data.data)[0];
			if (firstKey) {
				const friendlyMessage = `${firstKey}: ${err.data.data[firstKey].message}`;
				console.error(friendlyMessage);
				charErrorText = friendlyMessage;
				return;
			}
		}
	};

	onMount(async () => {
		if (!pb.authStore.isValid) {
			window.location.pathname = '/login';
			return;
		}
		record = await pb.collection('characters').getOne(page.params.id!);

		let rawTags = record.tag_names;

		try {
			// 如果它是字符串，尝试解析它
			if (typeof rawTags === 'string') {
				// 如果 rawTags 是 '"单相思"'，解析后得到 '单相思'
				// 如果 rawTags 是 '["A", "B"]'，解析后得到 ['A', 'B']
				const parsed = JSON.parse(rawTags);

				// 解析后判断结果：是数组直接赋值，是字符串则包装成数组
				if (Array.isArray(parsed)) {
					charTags = parsed;
				} else if (typeof parsed === 'string') {
					charTags = [parsed];
				}
			} else if (Array.isArray(rawTags)) {
				// 如果 SDK 已经自动帮你转成了数组
				charTags = rawTags;
			}
		} catch (e) {
			// 如果解析失败（比如 rawTags 本身就是个普通字符串 "单相思" 而不是 '"单相思"'）
			charTags = rawTags ? [rawTags] : [];
		}

		charName = record.name;
		charDescription = record.description;
		charOrigin = record.origin;
	});
</script>

<div class="flex min-h-[70vh] flex-col items-center justify-center px-4 py-12">
	<div class="card w-full max-w-md bg-base-100">
		{#if record}
			<div class="card-body flex flex-col gap-4 p-8">
				<h2 class="mb-2 text-center text-3xl font-bold">Update {record?.name}</h2>
				<p class="mb-6 text-center text-base-content/60">
					All fields must be filled out unless otherwise specified.
				</p>
				<div class="flex flex-col gap-4">
					<label class="form-control w-full">
						<div class="label"><span class="label-text font-medium">Name</span></div>
						<input
							type="text"
							name="characterName"
							placeholder="XXXX & YYYY"
							class="input-bordered input w-full"
							autocomplete="name"
							bind:value={charName}
							required
						/>
					</label>

					<label class="form-control w-full">
						<div class="label">
							<span class="label-text font-medium">Description</span>
						</div>
						<textarea
							name="characterpDescription"
							placeholder="description"
							class="textarea w-full"
							maxlength="1000"
							bind:value={charDescription}
							required
						></textarea>
					</label>

					<label class="form-control w-full">
						<div class="label"><span class="label-text font-medium">Origin</span></div>
						<input
							type="text"
							name="characterOrigin"
							placeholder="https://pdnode.com"
							class="input-bordered input w-full"
							autocomplete="name"
							bind:value={charOrigin}
						/>
					</label>

					<div class="form-control w-full max-w-xs">
						<label class="label" for="file-upload">
							<span class="label-text">Pictures (Option)</span>
						</label>

						<input
							type="file"
							id="file-upload"
							class="file-input-bordered file-input w-full file-input-primary"
							accept="image/*"
							multiple
							bind:files={charPictures}
						/>
					</div>

					<div class="flex w-full max-w-sm flex-col">
						<div class="label">
							<span class="label-text font-medium">Tags (Option)</span>
						</div>

						<input
							type="text"
							placeholder="Enter the tags and press Enter...."
							class="input-bordered input w-full"
							bind:value={charTagsCurrentInput}
							onkeydown={newCharAddTag}
						/>

						<div class="mt-2 flex flex-wrap gap-2">
							{#each charTags as tag, i}
								<div class="badge gap-2 badge-soft p-3 badge-primary">
									{tag}
									<button type="button" onclick={() => charTags.splice(i, 1)} class="text-xs"
										>✕</button
									>
								</div>
							{/each}
						</div>
					</div>

					<p class="text-red-500">{charErrorText}</p>

					<button type="submit" class="btn w-full text-lg btn-primary" onclick={updateCharacter}
						>Update</button
					>
				</div>
			</div>
		{:else}
			<div role="alert" class="alert alert-soft alert-error">
				<span>The record could not be found.</span>
			</div>
		{/if}
	</div>
</div>
