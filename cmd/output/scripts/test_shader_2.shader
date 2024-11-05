//

textures/testmap/test_shader_2 {
	qer_editorimage textures/testmap/test_shader_qer.tga
	surfaceparm nolightmap
	surfaceparm nonsolid
	surfaceparm trans
	qer_trans 0.75
	{
		map textures/testmap/test_shader_4.tga
		blendfunc add
		tcMod turb 1 0.75 -0.25 0.075
		tcMod scroll -0.05 -0.05
	}
	{
		map textures/testmap/test_shader_5.tga
		blendfunc gl_dst_color gl_one
		tcMod turb 0 0.5 0.25 0.075
		tcMod scroll 0.05 0.05
	}
}
