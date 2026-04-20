<script lang="ts">
	import pb from '$lib/pocketbase';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';

	let errorText = $state('');
	let tags = $state<string[]>([]);
	let currentInput = $state('');
	let cpName = $state('');
	let cpDescription = $state('');
	let cpPictures: FileList | undefined = $state();
	let cpCharacter1 = $state('');
	let cpCharacter2 = $state('');

	const updateCP = async () => {
		errorText = '';

		if (cpCharacter1 === cpCharacter2) {
			errorText = 'Character #1 and Character #2 cannot be the same';
			return;
		}

		const formData = new FormData();
		formData.append('name', cpName);
		formData.append('description', cpDescription);
		formData.append('owner', pb.authStore.record!.id);
		if (cpPictures && cpPictures.length > 0) {
			for (let file of cpPictures) {
				formData.append('images', file);
			}
		}

		tags.forEach((tag) => {
			formData.append('tag_names', tag);
		});
		formData.delete('characters');
		formData.append('characters', cpCharacter1);
		formData.append('characters', cpCharacter2);

		try {
			await pb.collection('cps').update(page.params.id!, formData);
			goto(`/cps/${page.params.id}`);
		} catch (err: any) {
			errorText = err.data.data?.message ?? 'Update failed. Please try again.';
			console.error(err.data);

			const firstKey = Object.keys(err.data.data)[0];
			if (firstKey) {
				const friendlyMessage = `${firstKey}: ${err.data.data[firstKey].message}`;
				console.error(friendlyMessage);
				errorText = friendlyMessage;
				return;
			}
		}
	};

	function addTag(e: KeyboardEvent) {
		if (e.key === 'Enter' && currentInput.trim()) {
			e.preventDefault();
			if (tags.length >= 10) {
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

	// Get Char(s)
	let characters = $state<any[]>([]);
	let record = $state<any>();
	const reloadCharacters = async () => {
		characters = await pb.collection('characters').getFullList({
			filter: `owner = "${pb.authStore.record!.id}"`,
			sort: '-created'
		});
	};
	onMount(async () => {
		if (!pb.authStore.isValid) {
			window.location.pathname = '/login';
			return;
		}
		reloadCharacters();

		record = await pb.collection('cps').getOne(page.params.id!);

		let rawTags = record.tag_names;

		try {
			// 如果它是字符串，尝试解析它
			if (typeof rawTags === 'string') {
				// 如果 rawTags 是 '"单相思"'，解析后得到 '单相思'
				// 如果 rawTags 是 '["A", "B"]'，解析后得到 ['A', 'B']
				const parsed = JSON.parse(rawTags);

				// 解析后判断结果：是数组直接赋值，是字符串则包装成数组
				if (Array.isArray(parsed)) {
					tags = parsed;
				} else if (typeof parsed === 'string') {
					tags = [parsed];
				}
			} else if (Array.isArray(rawTags)) {
				// 如果 SDK 已经自动帮你转成了数组
				tags = rawTags;
			}
		} catch (e) {
			// 如果解析失败（比如 rawTags 本身就是个普通字符串 "单相思" 而不是 '"单相思"'）
			tags = rawTags ? [rawTags] : [];
		}

		cpName = record.name;
		cpDescription = record.description;
		cpCharacter1 = record.characters[0];
		cpCharacter2 = record.characters[1];
	});
</script>

<div class="flex min-h-[70vh] flex-col items-center justify-center px-4 py-12">
	<div class="card w-full max-w-md bg-base-100">
		{#if record}
			<div class="card-body p-8">
				<h2 class="mb-2 text-center text-3xl font-bold">Update {record.name}</h2>
				<p class="mb-6 text-center text-base-content/60">
					All fields must be filled out unless otherwise specified.
				</p>
				<div class="flex flex-col gap-4">
					<label class="form-control w-full">
						<div class="label"><span class="label-text font-medium">Name</span></div>
						<input
							type="text"
							name="cpName"
							placeholder="XXXX & YYYY"
							class="input-bordered input w-full"
							autocomplete="name"
							bind:value={cpName}
							
						/>
					</label>

					<label class="form-control w-full">
						<div class="label"><span class="label-text font-medium">Description</span></div>
						<textarea
							name="cpDescription"
							placeholder="description"
							class="textarea w-full"
							maxlength="1000"
							bind:value={cpDescription}
						></textarea>
					</label>

					<div class="form-control w-full max-w-xs">
						<label class="label" for="file-upload">
							<span class="label-text">Pictures</span>
						</label>

						<input
							type="file"
							id="file-upload"
							class="file-input-bordered file-input w-full file-input-primary"
							accept="image/*"
							bind:files={cpPictures}
							multiple
							
						/>
					</div>

					<div class="flex w-full max-w-sm flex-col">
						<div class="label"><span class="label-text font-medium">Tags</span></div>

						<input
							type="text"
							placeholder="Enter the tags and press Enter...."
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
					</div>

					<label class="form-control w-full">
						<div class="label"><span class="label-text font-medium">Character #1</span></div>
						<select
							name="cpCharacter1"
							placeholder="Select a character"
							class="select-bordered select w-full"
							autocomplete="name"
							bind:value={cpCharacter1}
						>
							<option value="" selected disabled>Select a character</option>
							{#each characters as character}
								<option value={character.id}>{character.name}</option>
							{/each}
						</select>
					</label>

					<label class="form-control w-full">
						<div class="label"><span class="label-text font-medium">Character #2</span></div>
						<select
							name="cpCharacter2"
							placeholder="Select a character"
							class="select-bordered select w-full"
							autocomplete="name"
							bind:value={cpCharacter2}
						>
							<option value="" selected disabled>Select a character</option>
							{#each characters as character}
								<option value={character.id}>{character.name}</option>
							{/each}
						</select>
					</label>

					<p class="text-red-500">{errorText}</p>

					<div class="form-control mt-6">
						<button onclick={updateCP} type="submit" class="btn w-full text-lg btn-primary"
							>Update</button
						>
					</div>
				</div>
			</div>
		{:else}
			<div role="alert" class="alert alert-soft alert-error">
				<span>The record could not be found.</span>
			</div>
		{/if}
	</div>
</div>
