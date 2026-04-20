<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import pb from '$lib/pocketbase';
	import { onMount } from 'svelte';
    import TagInput from '$lib/components/ui/TagInput.svelte';
	import { toFormData } from '$lib/utils/api';
    

	let tags = $state<string[]>([]);
	let errorText = $state('');

	let record = $state<any>();

	let name = $state('');
	let description = $state('');
	let origin = $state('');
	let pictures: FileList | undefined = $state();

	const updateCharacter = async () => {
		errorText = '';

		try {
			await pb.collection('characters').update(page.params.id!, toFormData({
                name,
                description,
                origin,
                tag_names: tags,
                images: pictures,
                owner: record.owner
            }));
			goto(`/characters/${page.params.id}`);
		} catch (err: any) {
			errorText = err.data.data?.message ?? 'Update failed. Please try again.';

			const firstKey = Object.keys(err.data.data)[0];
			if (firstKey) {
				const friendlyMessage = `${firstKey}: ${err.data.data[firstKey].message}`;
				console.error(friendlyMessage);
				errorText = friendlyMessage;
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

		name = record.name;
		description = record.description;
		origin = record.origin;
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
							bind:value={name}
							required
						/>
					</label>

					<label class="form-control w-full">
						<div class="label">
							<span class="label-text font-medium">Description</span>
						</div>
						<textarea
							name="characterDescription"
							placeholder="description"
							class="textarea w-full"
							maxlength="1000"
							bind:value={description}
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
							bind:value={origin}
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
							bind:files={pictures}
						/>
					</div>

			        <TagInput title="Your Character Tags (Option)" bind:tags />

					<p class="text-red-500">{errorText}</p>

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
